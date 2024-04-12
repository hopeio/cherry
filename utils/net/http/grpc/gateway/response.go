package gateway

import (
	"context"
	httpi "github.com/hopeio/cherry/utils/net/http"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hopeio/cherry/protobuf/errorcode"
	"github.com/hopeio/cherry/protobuf/response"
	"github.com/hopeio/cherry/utils/encoding/protobuf/jsonpb"
	"github.com/hopeio/cherry/utils/net/http/grpc/reconn"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func ResponseHook(ctx context.Context, writer http.ResponseWriter, message proto.Message) error {
	if res, ok := message.(*response.HttpResponse); ok {
		hlen := len(res.Header)
		for i := 0; i < hlen && i+1 < hlen; i += 2 {
			writer.Header().Set(res.Header[i], res.Header[i+1])
		}
		writer.WriteHeader(int(res.StatusCode))
	}
	/*	if message == nil{
		*(&message) = &response.TinyRep{Message: "OK"}
	}*/
	return nil
}

func HTTPError(ctx *gin.Context, err error) {

	s, ok := status.FromError(err)
	if ok && s.Code() == 14 && strings.HasSuffix(s.Message(), `refused it."`) {
		//提供一个思路，这里应该是哪条连接失败重连哪条，不能这么粗暴，map的key是个关键
		if len(reconn.ReConnectMap) > 0 {
			for _, f := range reconn.ReConnectMap {
				f()
			}
		}
	}

	const fallback = `{"code": 14, "message": "failed to marshal error message"}`

	delete(ctx.Request.Header, httpi.HeaderTrailer)
	contentType := jsonpb.JsonPb.ContentType(nil)
	ctx.Header(httpi.HeaderContentType, contentType)

	se := &errorcode.ErrRep{Code: errorcode.ErrCode(s.Code()), Message: s.Message()}

	buf, merr := jsonpb.JsonPb.Marshal(se)
	if merr != nil {
		grpclog.Infof("Failed to marshal error message %q: %v", se, merr)
		ctx.Status(http.StatusInternalServerError)
		if _, err := io.WriteString(ctx.Writer, fallback); err != nil {
			grpclog.Infof("Failed to write response: %v", err)
		}
		return
	}

	if _, err := ctx.Writer.Write(buf); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}

}

type responseBody interface {
	XXX_ResponseBody() interface{}
}

func ForwardResponseMessage(ctx *gin.Context, md runtime.ServerMetadata, message proto.Message) {

	handleForwardResponseServerMetadata(ctx.Writer, md.HeaderMD)
	handleForwardResponseTrailerHeader(ctx.Writer, md.TrailerMD)

	if res, ok := message.(*response.HttpResponse); ok {
		hlen := len(res.Header)
		for i := 0; i < hlen && i+1 < hlen; i += 2 {
			ctx.Header(res.Header[i], res.Header[i+1])
		}
		ctx.Status(int(res.StatusCode))
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
	if rb, ok := message.(responseBody); ok {
		buf, err = jsonpb.JsonPb.Marshal(rb.XXX_ResponseBody())
	} else {
		buf, err = jsonpb.JsonPb.Marshal(message)
	}
	if err != nil {
		grpclog.Infof("Marshal error: %v", err)
		HTTPError(ctx, err)
		return
	}

	if _, err = ctx.Writer.Write(buf); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}

	handleForwardResponseTrailer(ctx.Writer, md.TrailerMD)
}
