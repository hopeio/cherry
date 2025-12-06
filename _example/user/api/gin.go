/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hopeio/cherry/_example/protobuf/user"
	"github.com/hopeio/cherry/_example/user/service"
)

func GinRegister(app *gin.Engine) {
	user.RegisterUserServiceHandlerServer(app, service.GetUserService())
	//oauth.RegisterOauthServiceHandlerServer(app, service.GetOauthService())
	app.StaticFS("/oauth/login", http.Dir("./static/login.html"))
	app.GET("/api/test/:id", service.Test)
}
