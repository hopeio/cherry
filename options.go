/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package cherry

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	"google.golang.org/grpc"
)

type Option func(server *Server)

func WithContext(ctx context.Context) Option {
	return func(server *Server) {
		server.BaseContext = ctx
	}
}

func WithHttp(http Http) Option {
	return func(server *Server) {
		server.Http = http
	}
}

func WithHTTP3(http3 Http3) Option {
	return func(server *Server) {
		server.HTTP3 = http3
	}
}

func WithGrpcHandler(handler func(*grpc.Server)) Option {
	return func(server *Server) {
		server.GrpcHandler = handler
	}
}

func WithGinHandler(handler func(*gin.Engine)) Option {
	return func(server *Server) {
		server.GinHandler = handler
	}
}

func WithGrpc(option GrpcConfig) Option {
	return func(server *Server) {
		server.Grpc = option
	}
}

func WithCors(cors cors.Options) Option {
	return func(server *Server) {
		server.Cors.Enabled = true
		server.Cors.Options = cors
	}
}

func WithTelemetry(telemetry TelemetryConfig) Option {
	return func(server *Server) {
		server.Telemetry = telemetry
	}
}
