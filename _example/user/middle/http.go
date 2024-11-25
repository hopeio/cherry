/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

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
