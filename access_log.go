/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package cherry

import (
	"github.com/hopeio/context/httpctx"
	"github.com/hopeio/utils/log"
	"go.uber.org/zap"
)

type AccessLogParam struct {
	Method, Url       string
	ReqBody, RespBody []byte
	Code              int
}

type AccessLog = func(ctxi *httpctx.Context, pram *AccessLogParam)

func defaultAccessLog(ctxi *httpctx.Context, param *AccessLogParam) {
	// log 里time now 浪费性能
	if ce := log.Default().Logger.Check(zap.InfoLevel, "access"); ce != nil {
		ce.Write(zap.String("url", param.Url),
			zap.String("method", param.Method),
			zap.ByteString("body", param.ReqBody),
			zap.String("traceId", ctxi.TraceID()),
			// 性能
			zap.Duration("processTime", ce.Time.Sub(ctxi.RequestAt.Time)),
			zap.ByteString("result", param.RespBody),
			zap.String("auth", ctxi.AuthInfoRaw),
			zap.Int("status", param.Code))
	}
}
