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
		} else if strings.HasPrefix(param.RequestRecorder.ContentType, httpx.ContentTypeProtobuf) {
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
			respBodyField = zap.Reflect("resp", json.RawMessage(param.ResponseRecorder.Raw))
		} else if strings.HasPrefix(param.ResponseRecorder.ContentType, httpx.ContentTypeProtobuf) {
			respBodyField = zap.String("resp", param.ResponseRecorder.Value.(fmt.Stringer).String())
		} else {
			respBodyField = zap.String("resp", stringsx.FromBytes(param.ResponseRecorder.Raw))
		}
	}
	authField := zap.Skip()
	metadata := ctx.Value(httpx.RequestMetadataKey).(*httpx.RequestMetadata)
	if metadata.Auth != nil && len(metadata.Auth.Raw) > 0 {
		zap.Reflect("auth", json.RawMessage(metadata.Auth.Raw))
	}

	if ce := log.NoCallerLogger().Logger.Check(zap.InfoLevel, "access"); ce != nil {
		ce.Write(zap.String("url", param.Url),
			zap.String("method", param.Method),
			reqBodyField,
			zap.String("traceId", metadata.TraceId),
			zap.Duration("duration", ce.Time.Sub(metadata.RequestAt)),
			respBodyField,
			authField,
			zap.Int("status", param.Code))
	}
}

type GrpcAccessLogParam struct {
	Method    string
	req, resp any
	err       error
}

type GrpcAccessLog = func(ctx context.Context, pram *GrpcAccessLogParam)

func DefaultGrpcAccessLog(ctx context.Context, param *GrpcAccessLogParam) {
	respBodyField := zap.Skip()
	codeField := zap.Int32("code", 0)
	if param.err != nil {
		s, _ := status.FromError(param.err)
		codeField = zap.Int32("code", int32(s.Code()))
		respBodyField = zap.String("resp", s.Message())
	} else {
		respBodyField = zap.String("resp", param.resp.(fmt.Stringer).String())
	}
	authField := zap.Skip()
	metadata := ctx.Value(httpx.RequestMetadataKey).(*httpx.RequestMetadata)
	if metadata.Auth != nil && len(metadata.Auth.Raw) > 0 {
		zap.Reflect("auth", json.RawMessage(metadata.Auth.Raw))
	}
	if ce := log.NoCallerLogger().Logger.Check(zap.InfoLevel, "access"); ce != nil {
		ce.Write(zap.String("url", param.Method),
			zap.String("method", "grpc"),
			zap.String("body", param.req.(fmt.Stringer).String()),
			zap.String("traceId", metadata.TraceId),
			zap.Duration("duration", ce.Time.Sub(metadata.RequestAt)),
			codeField,
			respBodyField,
			authField,
		)
	}
}
