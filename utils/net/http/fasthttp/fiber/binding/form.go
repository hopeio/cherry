package binding

import (
	"github.com/gofiber/fiber/v3"
	"github.com/hopeio/cherry/utils/net/http/binding"
	binding2 "github.com/hopeio/cherry/utils/net/http/fasthttp/binding"
)

type formBinding struct{}
type formPostBinding struct{}
type formMultipartBinding struct{}

func (formMultipartBinding) Name() string {
	return "multipart/form-data"
}

func (formMultipartBinding) Bind(ctx fiber.Ctx, obj interface{}) error {
	if err := binding.MappingByPtr(obj, (*binding2.MultipartRequest)(ctx.Request()), binding.Tag); err != nil {
		return err
	}

	return binding.Validate(obj)
}

func (formBinding) Name() string {
	return "form"
}

func (formBinding) Bind(ctx fiber.Ctx, obj interface{}) error {
	if err := binding.MapForm(obj, (*binding2.ArgsSource)(ctx.Request().PostArgs())); err != nil {
		return err
	}
	return binding.Validate(obj)
}

func (formPostBinding) Name() string {
	return "form-urlencoded"
}

func (formPostBinding) Bind(ctx fiber.Ctx, obj interface{}) error {
	if err := binding.MapForm(obj, (*binding2.ArgsSource)(ctx.Request().PostArgs())); err != nil {
		return err
	}
	return binding.Validate(obj)
}
