package binding

import (
	"github.com/hopeio/cherry/utils/net/http/request/binding"
	"github.com/valyala/fasthttp"
)

type xmlBinding struct{}

func (xmlBinding) Name() string {
	return "xml"
}

func (x xmlBinding) Bind(req *fasthttp.RequestCtx, obj interface{}) error {
	return binding.DecodeXml(req.Request.Body(), obj)
}
