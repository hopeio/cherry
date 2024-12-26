/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package cherry

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/utils/net/http/grpc/web"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

type IOption interface {
	apply(server *Server)
}

type Option func(server *Server)

func WithContext(ctx context.Context) Option {
	return func(server *Server) {
		server.BaseContext = ctx
	}
}

func WithHttpAddr(addr string) Option {
	return func(server *Server) {
		server.Http.Addr = addr
	}
}

func WithHttp3Addr(addr string) Option {
	return func(server *Server) {
		server.Http3.Addr = addr
	}
}

func WithName(name string) Option {
	return func(server *Server) {
		server.Name = name
	}
}

func WithStopTimeout(stopTimeout time.Duration) Option {
	return func(server *Server) {
		server.StopTimeout = stopTimeout
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

func WithGrpcWeb(options ...web.Option) Option {
	return func(server *Server) {
		server.EnableGrpcWeb = true
		server.GrpcWebOptions = append(server.GrpcWebOptions, options...)
	}
}

func WithOnStart(onStart func(context.Context)) Option {
	return func(server *Server) {
		server.OnStart = onStart
	}
}

func WithOnStop(onStop func(context.Context)) Option {
	return func(server *Server) {
		server.OnStop = onStop
	}
}

func WithGrpcOptions(options ...grpc.ServerOption) Option {
	return func(server *Server) {
		server.GrpcOptions = append(server.GrpcOptions, options...)
	}
}

func WithCors(cors *cors.Options) Option {
	return func(server *Server) {
		server.EnableCors = true
		server.Cors = cors
	}
}

func WithMiddlewares(middlewares ...http.HandlerFunc) Option {
	return func(server *Server) {
		server.Middlewares = append(server.Middlewares, middlewares...)
	}
}

func WithTelemetry(telemetry TelemetryConfig) Option {
	return func(server *Server) {
		server.EnableTelemetry = true
		server.TelemetryConfig = telemetry
	}
}
