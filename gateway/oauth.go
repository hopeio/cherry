/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package gateway

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hopeio/gox/reflect/mtos"
	"github.com/hopeio/protobuf/oauth"

	httpx "github.com/hopeio/gox/net/http"
	"github.com/hopeio/protobuf/response"
	"google.golang.org/grpc/metadata"
)

type OauthServiceServer interface {
	OauthAuthorize(context.Context, *oauth.OauthReq) (*response.HttpResponse, error)
	OauthToken(context.Context, *oauth.OauthReq) (*response.HttpResponse, error)
}

func RegisterOauthServiceHandlerServer(r *gin.Engine, server OauthServiceServer) {
	r.GET("/oauth/authorize", func(ctx *gin.Context) {
		var protoReq oauth.OauthReq
		mtos.DefaultDecoder().Decode(&protoReq, ctx.Request.URL.Query())
		res, _ := server.OauthAuthorize(
			metadata.NewIncomingContext(
				ctx.Request.Context(),
				metadata.MD{"auth": {httpx.GetToken(ctx.Request)}}),
			&protoReq)

		res.Respond(ctx, ctx.Writer)
	})

	r.POST("/oauth/access_token", func(ctx *gin.Context) {
		var protoReq oauth.OauthReq
		mtos.DefaultDecoder().Decode(&protoReq, ctx.Request.PostForm)
		res, _ := server.OauthToken(ctx.Request.Context(), &protoReq)
		res.Respond(ctx, ctx.Writer)
	})
}
