package gin

import (
	"github.com/gin-gonic/gin"
	httpi "github.com/hopeio/cherry/utils/net/http/debug"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Debug(r *gin.Engine) {
	r.Any("/debug/*path", Wrap(httpi.StackHandler()))
	// Register Metrics metrics handler.
	r.Any("/metrics", Wrap(promhttp.Handler()))
}
