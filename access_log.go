/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package cherry

import (
	"encoding/json"

	"github.com/hopeio/context/httpctx"
	"github.com/hopeio/gox/log"
	stringsx "github.com/hopeio/gox/strings"
	"go.uber.org/zap"
)

type Body struct {
	IsJson bool
	Data   []byte
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
		if param.ReqBody.IsJson {
			reqBodyField = zap.Reflect("body", json.RawMessage(param.ReqBody.Data))
		} else {
			reqBodyField = zap.String("body", stringsx.BytesToString(param.ReqBody.Data))
		}
	}
	respBodyField := zap.Skip()
	if len(param.RespBody.Data) > 0 {
		if param.RespBody.IsJson {
			respBodyField = zap.Reflect("result", json.RawMessage(param.RespBody.Data))
		} else {
			respBodyField = zap.String("result", stringsx.BytesToString(param.RespBody.Data))
		}
	}
	// log 里time now 浪费性能
	if ce := log.NoCallerLogger().Logger.Check(zap.InfoLevel, "access"); ce != nil {
		ce.Write(zap.String("url", param.Url),
			zap.String("method", param.Method),
			reqBodyField,
			zap.String("traceId", ctxi.TraceID()),
			// 性能
			zap.Duration("duration", ce.Time.Sub(ctxi.RequestAt.Time)),
			respBodyField,
			zap.String("auth", ctxi.AuthInfoRaw),
			zap.Int("status", param.StatusCode))
	}
}
