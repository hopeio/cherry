package middle

import (
	"github.com/gofiber/fiber/v3"
	"net/http"

	"github.com/hopeio/utils/log"
)

func Log(w http.ResponseWriter, r *http.Request) {
	log.Debug(r.RequestURI)
}

func FiberLog(ctx fiber.Ctx) error {
	log.Debug(ctx.BaseURL())
	return ctx.Next()
}
