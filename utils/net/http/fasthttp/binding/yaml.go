package binding

import (
	"github.com/hopeio/cherry/utils/net/http/binding"
	"github.com/valyala/fasthttp"
)

type yamlBinding struct{}

func (yamlBinding) Name() string {
	return "yaml"
}

func (y yamlBinding) Bind(req *fasthttp.RequestCtx, obj interface{}) error {
	return binding.DecodeYaml(req.Request.Body(), obj)
}
