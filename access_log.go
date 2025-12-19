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
	if len(param.Request.Raw) > 0 || param.Request.Value != nil || param.Request.Body != nil {
		if param.Request.Raw == nil && param.Request.Body != nil {
			param.Request.Raw = param.Request.Body.Bytes()
		}
		if strings.HasPrefix(param.Request.ContentType, httpx.ContentTypeJson) {
			reqBodyField = zap.Reflect("body", json.RawMessage(param.Request.Raw))
		} else if strings.HasPrefix(param.Request.ContentType, httpx.ContentTypeProtobuf) {
			reqBodyField = zap.String("body", param.Request.Value.(fmt.Stringer).String())
		} else {
			reqBodyField = zap.String("body", stringsx.BytesToString(param.Request.Raw))
		}
	}
	respBodyField := zap.Skip()
	if len(param.Reponse.Raw) > 0 || param.Reponse.Value != nil || param.Reponse.Body != nil {
		if param.Reponse.Raw == nil && param.Reponse.Body != nil {
			param.Reponse.Raw = param.Reponse.Body.Bytes()
		}
		if strings.HasPrefix(param.Reponse.ContentType, httpx.ContentTypeJson) {
			respBodyField = zap.Reflect("resp", json.RawMessage(param.Reponse.Raw))
		} else if strings.HasPrefix(param.Reponse.ContentType, httpx.ContentTypeProtobuf) {
			reqBodyField = zap.String("resp", param.Reponse.Value.(fmt.Stringer).String())
		} else {
			respBodyField = zap.String("resp", stringsx.BytesToString(param.Reponse.Raw))
		}
	}
	// log 里time now 浪费性能
	if ce := log.NoCallerLogger().Logger.Check(zap.InfoLevel, "access"); ce != nil {
		ce.Write(zap.String("url", param.Url),
			zap.String("method", param.Method),
			reqBodyField,
			zap.String("traceId", ctxi.TraceID()),
			zap.Duration("duration", ce.Time.Sub(ctxi.RequestTime.Time)),
			respBodyField,
			zap.String("auth", ctxi.AuthRaw),
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
