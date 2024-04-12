package binding

import (
	"github.com/gofiber/fiber/v3"
	fbinding "github.com/hopeio/cherry/utils/net/http/fasthttp/request/binding"
	"github.com/hopeio/cherry/utils/net/http/request/binding"
)

type queryBinding struct{}

func (queryBinding) Name() string {
	return "query"
}

func (queryBinding) Bind(ctx fiber.Ctx, obj interface{}) error {
	values := ctx.Request().URI().QueryArgs()
	if err := binding.MapForm(obj, (*fbinding.ArgsSource)(values)); err != nil {
		return err
	}
	return binding.Validate(obj)
}
