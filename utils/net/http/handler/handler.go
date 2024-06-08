package handler

import (
	"github.com/hopeio/cherry/utils/net/http/binding"
	"github.com/hopeio/cherry/utils/types"
	"net/http"
)

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
