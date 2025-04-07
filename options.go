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
	"net/http"
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

func WithGrpcHandler(grpcHandler func(*grpc.Server)) Option {
	return func(server *Server) {
		server.GrpcHandler = grpcHandler
	}
}

func WithGinHandler(ginHandler func(*gin.Engine)) Option {
	return func(server *Server) {
		server.GinHandler = ginHandler
	}
}

func WithGrpc(option GrpcConfig) Option {
	return func(server *Server) {
		server.Grpc = option
	}
}

func WithCors(cors cors.Options) Option {
	return func(server *Server) {
		server.Cors.Enable = true
		server.Cors.Options = cors
	}
}

func WithMiddlewares(middlewares ...http.HandlerFunc) Option {
	return func(server *Server) {
		server.Middlewares = append(server.Middlewares, middlewares...)
	}
}

func WithTelemetry(telemetry TelemetryConfig) Option {
	return func(server *Server) {
		server.Telemetry = telemetry
	}
}
