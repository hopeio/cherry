package cherry

import (
	"github.com/hopeio/context/httpctx"
	"github.com/hopeio/utils/log"
	"go.uber.org/zap"
)

type AccessLog = func(ctxi *httpctx.Context, uri, method, body, result string, code int)

func defaultAccessLog(ctxi *httpctx.Context, uri, method, body, result string, code int) {
	// log 里time now 浪费性能
	if ce := log.Default().Logger.Check(zap.InfoLevel, "access"); ce != nil {
		ce.Write(zap.String("uri", uri),
			zap.String("method", method),
			zap.String("body", body),
			zap.String("traceId", ctxi.TraceID()),
			// 性能
			zap.Duration("processTime", ce.Time.Sub(ctxi.RequestAt.Time)),
			zap.String("result", result),
			zap.String("auth", ctxi.AuthInfoRaw),
			zap.Int("status", code))
	}
}
