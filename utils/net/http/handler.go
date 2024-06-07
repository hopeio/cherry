package http

import (
	"encoding/json"
	"github.com/hopeio/cherry/protobuf/errorcode"
	"github.com/hopeio/cherry/utils/log"
	"github.com/hopeio/cherry/utils/net/http/binding"
	http_fs "github.com/hopeio/cherry/utils/net/http/fs"
	"github.com/hopeio/cherry/utils/types"
	"go.uber.org/zap"
	"io"
	"net/http"
	"reflect"
)

type Handlers []http.Handler

func (hs Handlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, handler := range hs {
		handler.ServeHTTP(w, r)
	}
}

type HandlerFuncs []http.HandlerFunc

func (hs HandlerFuncs) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, handler := range hs {
		handler(w, r)
	}
}

func (hs HandlerFuncs) HandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, handler := range hs {
			handler(w, r)
		}
	}
}

func (hs *HandlerFuncs) Add(handler http.HandlerFunc) {
	*hs = append(*hs, handler)
}

// TODO
func commonHandler[REQ, RES any](method types.GrpcServiceMethod[*REQ, *RES]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := new(REQ)
		err := binding.Bind(r, req)
		if err != nil {
			return
		}
	})
}

func ResWriteReflect(w http.ResponseWriter, traceId string, result []reflect.Value) {
	if !result[1].IsNil() {
		err := errorcode.ErrHandle(result[1].Interface())
		log.Errorw(err.Error(), zap.String(log.FieldTraceId, traceId))
		json.NewEncoder(w).Encode(err)
		return
	}
	if info, ok := result[0].Interface().(*http_fs.File); ok {
		header := w.Header()
		header.Set(HeaderContentType, ContentBinaryHeaderValue)
		header.Set(HeaderContentDisposition, "attachment;filename="+info.Name)
		io.Copy(w, info.File)
		if flusher, canFlush := w.(http.Flusher); canFlush {
			flusher.Flush()
		}
		info.File.Close()
		return
	}
	json.NewEncoder(w).Encode(ResAnyData{
		Message: "OK",
		Details: result[0].Interface(),
	})
}
