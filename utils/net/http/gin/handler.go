package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/hopeio/cherry/utils/definition/types"
)

// TODO:
func commonHandler[REQ, RES any](method types.GrpcServiceMethod[*REQ, *RES]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := new(REQ)
		err := ctx.Bind(req)
		if err != nil {
			return
		}
	}
}
