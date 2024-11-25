/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/hopeio/protobuf/response"
	"github.com/hopeio/utils/net/http/grpc"
	gin2 "github.com/hopeio/utils/net/http/grpc/gateway/gin"
	"google.golang.org/protobuf/proto"
)

func ForwardResponseMessage(ctx *gin.Context, md grpc.ServerMetadata, message proto.Message) {
	if res, ok := message.(*response.HttpResponse); ok {
		hlen := len(res.Header)
		for i := 0; i < hlen && i+1 < hlen; i += 2 {
			ctx.Header(res.Header[i], res.Header[i+1])
		}
		ctx.Status(int(res.StatusCode))
		ctx.Writer.Write(res.Body)
		return
	}

	gin2.ForwardResponseMessage(ctx, md, message)
}
