package pickrouter

import (
	contexti "github.com/hopeio/cherry/context"
	"github.com/hopeio/cherry/context/httpctx"
	"github.com/hopeio/cherry/protobuf/errorcode"
	"github.com/hopeio/cherry/utils/encoding/json"
	httpi "github.com/hopeio/cherry/utils/net/http"
	"github.com/hopeio/cherry/utils/net/http/binding"
	http_fs "github.com/hopeio/cherry/utils/net/http/fs"
	"io"
	"net/http"
	"reflect"
)

func commonHandler(w http.ResponseWriter, req *http.Request, handle *reflect.Value, ps *Params, tracing bool) {
	handleTyp := handle.Type()
	handleNumIn := handleTyp.NumIn()
	if handleNumIn != 0 {
		params := make([]reflect.Value, handleNumIn)
		ctxi, s := httpctx.ContextFromRequest(httpctx.RequestCtx{
			Request:  req,
			Response: w,
		}, tracing)
		if s != nil {
			defer s.End()
		}
		for i := 0; i < handleNumIn; i++ {
			if handleTyp.In(i).ConvertibleTo(HttpContextType) {
				params[i] = reflect.ValueOf(ctxi)
			} else {
				params[i] = reflect.New(handleTyp.In(i).Elem())
				if ps != nil || req.URL.RawQuery != "" {
					src := req.URL.Query()
					if ps != nil {
						pathParam := *ps
						if len(pathParam) > 0 {
							for i := range pathParam {
								src.Set(pathParam[i].Key, pathParam[i].Value)
							}
						}
					}
					binding.Decode(params[i], src)
				}
				if req.Method != http.MethodGet {
					json.NewDecoder(req.Body).Decode(params[i].Interface())
				}
			}
		}
		result := handle.Call(params)
		ResHandler(ctxi, w, result)
	}
}

func ResHandler[T any](c *contexti.RequestContext[T], w http.ResponseWriter, result []reflect.Value) {
	if !result[1].IsNil() {
		err := errorcode.ErrHandle(result[1].Interface())
		c.HandleError(err)
		json.NewEncoder(w).Encode(err)
		return
	}
	if info, ok := result[0].Interface().(*http_fs.File); ok {
		header := w.Header()
		header.Set(httpi.HeaderContentType, httpi.ContentBinaryHeaderValue)
		header.Set(httpi.HeaderContentDisposition, "attachment;filename="+info.Name)
		io.Copy(w, info.File)
		if flusher, canFlush := w.(http.Flusher); canFlush {
			flusher.Flush()
		}
		info.File.Close()
		return
	}
	json.NewEncoder(w).Encode(httpi.ResAnyData{
		Message: "OK",
		Details: result[0].Interface(),
	})
}
