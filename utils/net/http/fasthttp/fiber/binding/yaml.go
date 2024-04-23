package binding

import (
	"github.com/gofiber/fiber/v3"
	"github.com/hopeio/cherry/utils/net/http/binding"
)

type yamlBinding struct{}

func (yamlBinding) Name() string {
	return "yaml"
}

func (y yamlBinding) Bind(ctx fiber.Ctx, obj interface{}) error {
	return binding.DecodeYaml(ctx.Body(), obj)
}
