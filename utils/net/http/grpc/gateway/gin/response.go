package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/hopeio/cherry/protobuf/response"
	"github.com/hopeio/cherry/utils/encoding/protobuf/jsonpb"
	httpi "github.com/hopeio/cherry/utils/net/http"
	"github.com/hopeio/cherry/utils/net/http/grpc"
	"github.com/hopeio/cherry/utils/net/http/grpc/gateway"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/proto"
)

func ForwardResponseMessage(ctx *gin.Context, md grpc.ServerMetadata, message proto.Message) {

	gateway.HandleForwardResponseServerMetadata(ctx.Writer, md.HeaderMD)
	gateway.HandleForwardResponseTrailerHeader(ctx.Writer, md.TrailerMD)

	if res, ok := message.(*response.HttpResponse); ok {
		hlen := len(res.Header)
		for i := 0; i < hlen && i+1 < hlen; i += 2 {
			ctx.Header(res.Header[i], res.Header[i+1])
		}
		ctx.Status(int(res.Status))
		ctx.Writer.Write(res.Body)
		return
	}

	contentType := jsonpb.JsonPb.ContentType(message)
	ctx.Header(httpi.HeaderContentType, contentType)

	if !message.ProtoReflect().IsValid() {
		ctx.Writer.Write(httpi.ResponseOk)
		return
	}

	var buf []byte
	var err error
	switch rb := message.(type) {
	case responseBody:
		buf, err = jsonpb.JsonPb.Marshal(rb.ResponseBody())
	case xxxResponseBody:
		buf, err = jsonpb.JsonPb.Marshal(rb.XXX_ResponseBody())
	default:
		buf, err = jsonpb.JsonPb.Marshal(message)
	}

	if err != nil {
		grpclog.Infof("Marshal error: %v", err)
		HttpError(ctx, err)
		return
	}

	if _, err = ctx.Writer.Write(buf); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}

	gateway.HandleForwardResponseTrailer(ctx.Writer, md.TrailerMD)
}

type xxxResponseBody interface {
	XXX_ResponseBody() interface{}
}

type responseBody interface {
	ResponseBody() interface{}
}
