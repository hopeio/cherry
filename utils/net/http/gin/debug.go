package gin

import (
	"github.com/gin-gonic/gin"
	httpi "github.com/hopeio/cherry/utils/net/http/handlers/debug"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Debug(r *gin.Engine) {
	r.Any("/debug/*path", Wrap(httpi.DebugHandler()))
	// Register Metrics metrics handler.
	r.Any("/metrics", Wrap(promhttp.Handler()))
}
