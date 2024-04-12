package response

import (
	"google.golang.org/protobuf/types/known/wrapperspb"
	"net/http"

	"google.golang.org/protobuf/proto"
)

type GoReply struct {
	Code    uint32
	Message string
	Details proto.Message
}

func (x *HttpResponse) GetContentType() string {
	hlen := len(x.Header)
	for i := 0; i < hlen && i+1 < hlen; i += 2 {
		if x.Header[i] == "Content-Type" {
			return x.Header[i+1]
		}
	}
	return ""
}

func (x *HttpResponse) Response(w http.ResponseWriter) {
	//我也是头一次知道要按顺序来的 response.wroteHeader
	//先设置请求头，再设置状态码，再写body
	//原因是http里每次操作都要判断wroteHeader(表示已经写过header了，不可以再写了)
	hlen := len(x.Header)
	for i := 0; i < hlen && i+1 < hlen; i += 2 {
		w.Header().Set(x.Header[i], x.Header[i+1])
	}
	w.WriteHeader(int(x.StatusCode))
	w.Write(x.Body)
}

var ResponseOk = &TinyRep{Message: "OK"}

type StringValue = wrapperspb.StringValue
type StringValueInput = wrapperspb.StringValue
