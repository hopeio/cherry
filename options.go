/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package cherry

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	httpx "github.com/hopeio/gox/net/http"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
)

type Option func(server *Server)

func WithContext(ctx context.Context) Option {
	return func(server *Server) {
		server.BaseContext = ctx
	}
}

func WithHttp(handler func(s *http.Server)) Option {
	return func(server *Server) {
		handler(&server.Server)
	}
}

func WithHttp2(handler func(s *http2.Server)) Option {
	return func(server *Server) {
		handler(&server.HTTP2)
	}
}

func WithHTTP3(handler func(s *Http3)) Option {
	return func(server *Server) {
		handler(&server.HTTP3)
	}
}

func WithGrpcHandler(handler func(*grpc.Server)) Option {
	return func(server *Server) {
		server.GrpcHandler = handler
	}
}

func WithGinHandler(handler func(*gin.Engine)) Option {
	return func(server *Server) {
		handler(server.GinServer)
	}
}

func WithGrpc(handler func(option *GrpcConfig)) Option {
	return func(server *Server) {
		handler(&server.Grpc)
	}
}

func WithCors(handler func(cors *cors.Options)) Option {
	return func(server *Server) {
		handler(&server.Cors.Options)
	}
}

func WithTelemetry(handler func(telemetry *TelemetryConfig)) Option {
	return func(server *Server) {
		handler(&server.Telemetry)
	}
}

func WithMiddleware(mw ...httpx.Middleware) Option {
	return func(server *Server) {
		server.Middlewares = append(server.Middlewares, mw...)
	}
}
