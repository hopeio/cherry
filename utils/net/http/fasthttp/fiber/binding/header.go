package binding

import (
	"github.com/gofiber/fiber/v3"
	fbinding "github.com/hopeio/cherry/utils/net/http/fasthttp/request/binding"
	"github.com/hopeio/cherry/utils/net/http/request/binding"
)

type headerBinding struct{}

func (headerBinding) Name() string {
	return "header"
}

func (headerBinding) Bind(req fiber.Ctx, obj interface{}) error {

	if err := binding.MappingByPtr(obj, (*fbinding.HeaderSource)(&req.Request().Header), binding.Tag); err != nil {
		return err
	}

	return binding.Validate(obj)
}
