package response

import (
	"errors"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io"
	"net/http"

	"google.golang.org/protobuf/proto"
)

type Reply struct {
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

func (x HttpResponse) MarshalProto(w io.Writer) {
	w.Write(x.Body)
}

func (x *HttpResponse) MarshalGQL(w io.Writer) {
	w.Write(x.Body)
}

func (x *HttpResponse) UnmarshalGQL(v interface{}) error {
	if i, ok := v.([]byte); ok {
		x.Body = i
		return nil
	}
	return errors.New("error type")
}

var ResponseOk = &TinyRep{Message: "OK"}

type StringValue = wrapperspb.StringValue

// graphql compatible
type StringValueInput = wrapperspb.StringValue

/*// graphql compatible
type HeaderEntry struct {
	Key   string
	Value string
}

type HttpResponseResolver struct {
}

// graphql compatible
func (receiver *HttpResponseResolver) Header(ctx context.Context, obj *HttpResponse) ([]*HeaderEntry, error) {
	var header []*HeaderEntry
	for k, v := range obj.Header {
		header = append(header, &HeaderEntry{Key: k, Value: v})
	}
	return header, nil
}
*/
