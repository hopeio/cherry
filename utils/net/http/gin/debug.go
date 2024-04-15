package gini

import (
	"github.com/gin-gonic/gin"
	httpi "github.com/hopeio/cherry/utils/net/http/debug"
	"github.com/hopeio/cherry/utils/net/http/gin/handler"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Debug(r *gin.Engine) {
	r.Any("/debug/*path", handler.Wrap(httpi.Debug()))
	// Register Metrics metrics handler.
	r.Any("/metrics", handler.Wrap(promhttp.Handler()))
}
