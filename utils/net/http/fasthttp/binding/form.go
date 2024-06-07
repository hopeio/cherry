package binding

import (
	"github.com/hopeio/cherry/utils/net/http/binding"
	"github.com/valyala/fasthttp"
)

type formPostBinding struct{}
type formMultipartBinding struct{}

func (formMultipartBinding) Name() string {
	return "multipart/form-data"
}

func (formMultipartBinding) Bind(req *fasthttp.RequestCtx, obj interface{}) error {
	if err := binding.MappingByPtr(obj, (*MultipartRequest)(&req.Request), binding.Tag); err != nil {
		return err
	}

	return binding.Validate(obj)
}

func (formPostBinding) Name() string {
	return "form-urlencoded"
}
func (formPostBinding) Bind(req *fasthttp.RequestCtx, obj interface{}) error {
	if err := binding.MapForm(obj, (*ArgsSource)(req.PostArgs())); err != nil {
		return err
	}
	return binding.Validate(obj)
}
