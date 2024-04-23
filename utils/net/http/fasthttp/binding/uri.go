package binding

import (
	"github.com/hopeio/cherry/utils/net/http/binding"
	stringsi "github.com/hopeio/cherry/utils/strings"
	"github.com/valyala/fasthttp"
	"reflect"
)

// fasthttp不支持uri路径参数
type uriBinding struct{}

func (uriBinding) Name() string {
	return "uri"
}

func (uriBinding) Bind(c *fasthttp.RequestCtx, obj interface{}) error {

	if err := binding.MappingByPtr(obj, (*uriSource)(c.URI()), binding.Tag); err != nil {
		return err
	}
	return binding.Validate(obj)
}

type uriSource fasthttp.URI

// TrySet tries to set a value by request's form source (like map[string][]string)
func (form *uriSource) TrySet(value reflect.Value, field reflect.StructField, tagValue string, opt binding.SetOptions) (isSet bool, err error) {
	return binding.SetByKV(value, field, form, tagValue, opt)
}

func (form *uriSource) Peek(key string) ([]string, bool) {
	v := stringsi.BytesToString((*fasthttp.URI)(form).LastPathSegment())
	return []string{v}, v != ""
}
