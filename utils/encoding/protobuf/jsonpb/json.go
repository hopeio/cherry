package jsonpb

import (
	"github.com/hopeio/cherry/utils/encoding/json"
	httpi "github.com/hopeio/cherry/utils/net/http"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hopeio/cherry/protobuf/response"
)

var JsonPb = &JSONPb{}

type JSONPb struct {
}

func (*JSONPb) ContentType(_ interface{}) string {
	return "application/json"
}

func (j *JSONPb) Marshal(v any) ([]byte, error) {
	if _, ok := v.(error); ok {
		return json.Marshal(v)
	}
	if msg, ok := v.(*wrapperspb.StringValue); ok {
		v = msg.Value
	}
	return json.Marshal(&httpi.ResAnyData{
		Code:    0,
		Message: "OK",
		Details: v,
	})
}

func (j *JSONPb) Name() string {
	return "jsonpb"
}

func (j *JSONPb) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (j *JSONPb) Delimiter() []byte {
	return []byte("\n")
}

// NewDecoder returns a runtime.Decoder which reads JSON stream from "r".
func (j *JSONPb) NewDecoder(r io.Reader) runtime.Decoder {
	return json.NewDecoder(r)
}

// NewEncoder returns an Encoder which writes JSON stream into "w".
func (j *JSONPb) NewEncoder(w io.Writer) runtime.Encoder {
	return json.NewEncoder(w)
}

func (j *JSONPb) ContentTypeFromMessage(v interface{}) string {
	if httpBody, ok := v.(*response.HttpResponse); ok {
		return httpBody.GetContentType()
	}
	return j.ContentType(v)
}
