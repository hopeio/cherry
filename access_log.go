/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package cherry

import (
	"encoding/json"
	"github.com/hopeio/context/httpctx"
	"github.com/hopeio/utils/log"
	stringsi "github.com/hopeio/utils/strings"
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
	if param.ReqBody.IsJson {
		reqBodyField = zap.Reflect("body", json.RawMessage(param.ReqBody.Data))
	} else {
		reqBodyField = zap.String("body", stringsi.BytesToString(param.ReqBody.Data))
	}
	respBodyField := zap.Skip()
	if param.RespBody.IsJson {
		respBodyField = zap.Reflect("body", json.RawMessage(param.RespBody.Data))
	} else {
		respBodyField = zap.String("body", stringsi.BytesToString(param.RespBody.Data))
	}
	// log 里time now 浪费性能
	if ce := log.Default().Logger.Check(zap.InfoLevel, "access"); ce != nil {
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
