/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package cherry

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hopeio/gox/log"
	httpx "github.com/hopeio/gox/net/http"
	stringsx "github.com/hopeio/gox/strings"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

type Body struct {
	ContentType string
	Raw         []byte
	Data        any
}

type AccessLogParam struct {
	Method, Url string
	*httpx.Recorder
	Metadata *Metadata
}

type AccessLog = func(ctx context.Context, pram *AccessLogParam)

func DefaultAccessLog(ctx context.Context, param *AccessLogParam) {
	reqBodyField := zap.Skip()
	if len(param.RequestRecorder.Raw) > 0 || param.RequestRecorder.Value != nil || param.RequestRecorder.Body != nil {
		if param.RequestRecorder.Raw == nil && param.RequestRecorder.Body != nil {
			param.RequestRecorder.Raw = param.RequestRecorder.Body.Bytes()
		}
		if strings.HasPrefix(param.RequestRecorder.ContentType, httpx.ContentTypeJson) {
			reqBodyField = zap.Reflect("body", json.RawMessage(param.RequestRecorder.Raw))
		} else if strings.HasPrefix(param.RequestRecorder.ContentType, httpx.ContentTypeXProtobuf) {
			reqBodyField = zap.String("body", param.RequestRecorder.Value.(fmt.Stringer).String())
		} else {
			reqBodyField = zap.String("body", stringsx.FromBytes(param.RequestRecorder.Raw))
		}
	}
	respBodyField := zap.Skip()
	if len(param.ResponseRecorder.Raw) > 0 || param.ResponseRecorder.Value != nil || param.ResponseRecorder.Body != nil {
		if param.ResponseRecorder.Raw == nil && param.ResponseRecorder.Body != nil {
			param.ResponseRecorder.Raw = param.ResponseRecorder.Body.Bytes()
		}
		if strings.HasPrefix(param.ResponseRecorder.ContentType, httpx.ContentTypeJson) {
			respBodyField = zap.Reflect("response", json.RawMessage(param.ResponseRecorder.Raw))
		} else if strings.HasPrefix(param.ResponseRecorder.ContentType, httpx.ContentTypeXProtobuf) {
			respBodyField = zap.String("response", param.ResponseRecorder.Value.(fmt.Stringer).String())
		} else {
			respBodyField = zap.String("response", stringsx.FromBytes(param.ResponseRecorder.Raw))
		}
	}

	if ce := log.NoCallerLogger().Logger.Check(zap.InfoLevel, "access"); ce != nil {
		ce.Write(zap.Inline(zap.DictObject(param.Metadata.AccessLogFields...)),
			zap.String("url", param.Url),
			zap.String("method", param.Method),
			reqBodyField,
			log.Context(ctx),
			zap.Duration("duration", ce.Time.Sub(param.Metadata.RequestAt)),
			respBodyField,
			zap.Int("status", param.StatusCode))
	}
}

type GrpcAccessLogParam struct {
	Method            string
	Request, Response any
	Err               error
	Metadata          *Metadata
}

type GrpcAccessLog = func(ctx context.Context, pram *GrpcAccessLogParam)

func DefaultGrpcAccessLog(ctx context.Context, param *GrpcAccessLogParam) {
	respBodyField := zap.Skip()
	codeField := zap.Int32("code", 0)
	if param.Err != nil {
		s, _ := status.FromError(param.Err)
		codeField = zap.Int32("code", int32(s.Code()))
		respBodyField = zap.String("response", s.Message())
	} else {
		respBodyField = zap.String("response", param.Response.(fmt.Stringer).String())
	}

	if ce := log.NoCallerLogger().Logger.Check(zap.InfoLevel, "access"); ce != nil {
		ce.Write(zap.Inline(zap.DictObject(param.Metadata.AccessLogFields...)),
			zap.String("url", param.Method),
			zap.String("method", "grpc"),
			zap.String("body", param.Request.(fmt.Stringer).String()),
			log.Context(ctx),
			zap.Duration("duration", ce.Time.Sub(param.Metadata.RequestAt)),
			codeField,
			respBodyField)
	}
}
