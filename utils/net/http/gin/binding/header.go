package binding

import (
	"github.com/hopeio/cherry/utils/net/http/binding"
	"net/http"
)

type headerBinding struct{}

func (headerBinding) Name() string {
	return "header"
}

func (headerBinding) Bind(req *http.Request, obj interface{}) error {

	if err := mapHeader(obj, req.Header); err != nil {
		return err
	}

	return Validate(obj)
}

func mapHeader(ptr interface{}, h map[string][]string) error {
	return binding.MappingByPtr(ptr, binding.HeaderSource(h), Tag)
}
