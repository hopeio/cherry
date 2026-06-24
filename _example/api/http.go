/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package api

import (
	"net/http"

	"github.com/hopeio/cherry/_example/service"
	"github.com/hopeio/gox/net/http/grpc/gateway"
)

func HttpHandler() http.Handler {
	mux:=http.NewServeMux()
	mux.HandleFunc("/api/signup", gateway.UnaryCall(service.GetUserService().Signup))
	mux.HandleFunc("/api/getUser", gateway.UnaryCall(service.GetUserService().GetUser))
	return mux
}
