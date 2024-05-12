package server

import (
	"github.com/hopeio/cherry/context/httpctx"
	"github.com/hopeio/cherry/utils/log"
	"go.uber.org/zap"
)

type AccessLog = func(ctxi *httpctx.Context, uri, method, body, result string, code int)

var defaultAccessLog = func(ctxi *httpctx.Context, uri, method, body, result string, code int) {
	// log 里time now 浪费性能
	if ce := log.Default().Logger.Check(zap.InfoLevel, "access"); ce != nil {
		ce.Write(zap.String("uri", uri),
			zap.String("method", method),
			zap.String("body", body),
			zap.String("traceId", ctxi.TraceID),
			// 性能
			zap.Duration("processTime", ce.Time.Sub(ctxi.RequestAt.Time)),
			zap.String("result", result),
			zap.String("auth", ctxi.AuthInfoRaw),
			zap.Int("status", code))
	}
}

func SetAccessLog(accessLog AccessLog) {
	if accessLog != nil {
		defaultAccessLog = accessLog
	}
}
