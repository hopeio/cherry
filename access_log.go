/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package cherry

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hopeio/gox/context/httpctx"
	"github.com/hopeio/gox/log"
	httpx "github.com/hopeio/gox/net/http"
	stringsx "github.com/hopeio/gox/strings"
	"github.com/hopeio/protobuf/response"
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

type AccessLog = func(ctxi *httpctx.Context, pram *AccessLogParam)

func DefaultAccessLog(ctxi *httpctx.Context, param *AccessLogParam) {
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
	if len(ctxi.AuthRaw) > 0 {
		zap.Reflect("auth", json.RawMessage(ctxi.AuthRaw))
	}
	// log 里time now 浪费性能
	if ce := log.NoCallerLogger().Logger.Check(zap.InfoLevel, "access"); ce != nil {
		ce.Write(zap.String("url", param.Url),
			zap.String("method", param.Method),
			reqBodyField,
			zap.String("traceId", ctxi.TraceID()),
			zap.Duration("duration", ce.Time.Sub(ctxi.RequestTime.Time)),
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

type GrpcAccessLog = func(ctxi *httpctx.Context, pram *GrpcAccessLogParam)

func DefaultGrpcAccessLog(ctxi *httpctx.Context, param *GrpcAccessLogParam) {
	respBodyField := zap.Skip()
	if param.err != nil {
		s, _ := status.FromError(param.err)
		se := &response.ErrResp{Code: int32(s.Code()), Msg: s.Message()}
		respBodyField = zap.String("resp", se.String())
	} else {
		respBodyField = zap.String("resp", param.resp.(fmt.Stringer).String())
	}

	if ce := log.NoCallerLogger().Logger.Check(zap.InfoLevel, "access"); ce != nil {
		ce.Write(zap.String("url", param.Method),
			zap.String("method", "grpc"),
			zap.String("body", param.req.(fmt.Stringer).String()),
			zap.String("traceId", ctxi.TraceID()),
			zap.Duration("duration", ce.Time.Sub(ctxi.RequestTime.Time)),
			respBodyField,
			zap.String("auth", ctxi.AuthRaw),
		)
	}
}
