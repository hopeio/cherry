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
	Data        []byte
}

type AccessLogParam struct {
	Method, Url       string
	ReqBody, RespBody Body
	StatusCode        int
}

type AccessLog = func(ctxi *httpctx.Context, pram *AccessLogParam)

func DefaultAccessLog(ctxi *httpctx.Context, param *AccessLogParam) {
	reqBodyField := zap.Skip()
	if len(param.ReqBody.Data) > 0 {
		if strings.HasPrefix(param.ReqBody.ContentType, httpx.ContentTypeJson) {
			reqBodyField = zap.Reflect("body", json.RawMessage(param.ReqBody.Data))
		} else {
			reqBodyField = zap.String("body", stringsx.BytesToString(param.ReqBody.Data))
		}
	}
	respBodyField := zap.Skip()
	if len(param.RespBody.Data) > 0 {
		if strings.HasPrefix(param.RespBody.ContentType, httpx.ContentTypeJson) {
			respBodyField = zap.Reflect("resp", json.RawMessage(param.RespBody.Data))
		} else {
			respBodyField = zap.String("resp", stringsx.BytesToString(param.RespBody.Data))
		}
	}
	// log 里time now 浪费性能
	if ce := log.NoCallerLogger().Logger.Check(zap.InfoLevel, "access"); ce != nil {
		ce.Write(zap.String("url", param.Url),
			zap.String("method", param.Method),
			reqBodyField,
			zap.String("traceId", ctxi.TraceID()),
			zap.Duration("duration", ce.Time.Sub(ctxi.RequestAt.Time)),
			respBodyField,
			zap.String("auth", ctxi.AuthInfoRaw),
			zap.Int("status", param.StatusCode))
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
		se := &response.CommonResp{Code: int32(s.Code()), Msg: s.Message()}
		respBodyField = zap.String("resp", se.String())
	} else {
		respBodyField = zap.String("resp", param.resp.(fmt.Stringer).String())
	}

	if ce := log.NoCallerLogger().Logger.Check(zap.InfoLevel, "access"); ce != nil {
		ce.Write(zap.String("url", param.Method),
			zap.String("method", "grpc"),
			zap.String("body", param.req.(fmt.Stringer).String()),
			zap.String("traceId", ctxi.TraceID()),
			zap.Duration("duration", ce.Time.Sub(ctxi.RequestAt.Time)),
			respBodyField,
			zap.String("auth", ctxi.AuthInfoRaw),
		)
	}
}
