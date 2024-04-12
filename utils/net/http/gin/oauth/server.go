package oauth

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/cherry/protobuf/oauth"
	"github.com/hopeio/cherry/protobuf/response"
	httpi "github.com/hopeio/cherry/utils/net/http/request"
	"github.com/hopeio/cherry/utils/net/http/request/binding"

	"google.golang.org/grpc/metadata"
)

type OauthServiceServer interface {
	OauthAuthorize(context.Context, *oauth.OauthReq) (*response.HttpResponse, error)
	OauthToken(context.Context, *oauth.OauthReq) (*response.HttpResponse, error)
}

func RegisterOauthServiceHandlerServer(r *gin.Engine, server OauthServiceServer) {
	r.GET("/oauth/authorize", func(ctx *gin.Context) {
		var protoReq oauth.OauthReq
		binding.DefaultDecoder().Decode(&protoReq, ctx.Request.URL.Query())
		res, _ := server.OauthAuthorize(
			metadata.NewIncomingContext(
				ctx.Request.Context(),
				metadata.MD{"auth": {httpi.GetToken(ctx.Request)}}),
			&protoReq)

		res.Response(ctx.Writer)
	})

	r.POST("/oauth/access_token", func(ctx *gin.Context) {
		var protoReq oauth.OauthReq
		binding.DefaultDecoder().Decode(&protoReq, ctx.Request.PostForm)
		res, _ := server.OauthToken(ctx.Request.Context(), &protoReq)
		res.Response(ctx.Writer)
	})
}
