package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/hopeio/cherry/protobuf/errcode"
	httpi "github.com/hopeio/cherry/utils/net/http"
)

func RespErrcode(ctx *gin.Context, code errcode.ErrCode) {
	httpi.RespErrcode(ctx.Writer, code)
}

func RespErr(ctx *gin.Context, err error) {
	httpi.RespErr(ctx.Writer, err)
}

func RespErrMsg(ctx *gin.Context, msg string) {
	httpi.RespErrMsg(ctx.Writer, msg)
}

func RespErrRep(ctx *gin.Context, rep *errcode.ErrRep) {
	httpi.RespErrRep(ctx.Writer, rep)
}

func Response(ctx *gin.Context, code errcode.ErrCode, msg string, data interface{}) {
	httpi.Response(ctx.Writer, code, msg, data)
}
