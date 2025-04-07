/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package cherry

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/utils/crypto/tls"
	"github.com/hopeio/utils/log"
	gini "github.com/hopeio/utils/net/http/gin"
	"github.com/hopeio/utils/net/http/grpc/web"
	"github.com/hopeio/utils/validation/validator"
	"github.com/quic-go/quic-go/http3"
	"github.com/rs/cors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

func NewServer(options ...Option) *Server {
	c := &Server{}
	c.Http.Addr = ":8080"
	c.Http.ReadTimeout = 5 * time.Second
	c.Http.WriteTimeout = 5 * time.Second
	gin.SetMode(gin.ReleaseMode)
	gin.DisableBindValidation()
	validator.DefaultValidator = nil // 自己做校验
	c.Cors.Enable = true
	c.Telemetry.Enable = true
	c.EnableDebugApi = true
	for _, option := range options {
		option(c)
	}
	return c
}

type Http struct {
	http.Server
	CertFile string
	KeyFile  string
	HttpOption
}

type Http3 struct {
	Enable bool
	http3.Server
	CertFile string
	KeyFile  string
}

type HttpOption struct {
	AccessLog          AccessLog
	ExcludeLogPrefixes []string
	IncludeLogPrefixes []string
	StaticFs           []StaticFsConfig
	Middlewares        []http.HandlerFunc
}

type StaticFsConfig struct {
	Prefix string
	Root   string
}

type Server struct {
	Http
	HTTP2          http2.Server
	HTTP3          Http3
	Gin            gini.Config
	Cors           CorsConfig
	Grpc           GrpcConfig
	EnableDebugApi bool
	ApiDoc         ApiDocConfig
	Telemetry      TelemetryConfig
	BaseContext    context.Context
	// 注册 grpc 服务
	GrpcHandler func(*grpc.Server)
	// 注册 gin 服务
	GinHandler func(*gin.Engine)
}

type ApiDocConfig struct {
	Enable         bool
	UriPrefix, Dir string
}

type GrpcConfig struct {
	EnableGrpcWeb  bool
	GrpcWebOptions []web.Option
	Options        []grpc.ServerOption
}

type CorsConfig struct {
	Enable bool
	cors.Options
}

type TelemetryConfig struct {
	Enable         bool
	EnableMetrics  bool
	EnableTracing  bool
	otelhttpOpts   []otelhttp.Option
	otelgrpcOpts   []otelgrpc.Option
	propagator     propagation.TextMapPropagator
	tracerProvider *sdktrace.TracerProvider
	meterProvider  *sdkmetric.MeterProvider
}

func (c *TelemetryConfig) SetOtelhttpHandlerOpts(otelhttpOpts []otelhttp.Option) {
	c.otelhttpOpts = otelhttpOpts
}

func (c *TelemetryConfig) SetOtelgrpcOptsHandlerOpts(otelgrpcOpts []otelgrpc.Option) {
	c.otelgrpcOpts = otelgrpcOpts
}

func (c *TelemetryConfig) SetTextMapPropagator(propagator propagation.TextMapPropagator) {
	c.propagator = propagator
}

func (c *TelemetryConfig) SetTracerProvider(tracerProvider *sdktrace.TracerProvider) {
	c.tracerProvider = tracerProvider
}

func (c *TelemetryConfig) SetMeterProvider(meterProvider *sdkmetric.MeterProvider) {
	c.meterProvider = meterProvider
}

func (s *Server) Init() {
	if s.Http.Addr == "" {
		s.Http.Addr = ":8080"
	}
	if s.HTTP3.Enable && s.HTTP3.Addr == "" {
		s.HTTP3.Addr = ":8080"
	}
	log.DurationNotify("ReadTimeout", s.Http.ReadTimeout, time.Second)
	log.DurationNotify("WriteTimeout", s.Http.WriteTimeout, time.Second)
	if s.Http.CertFile != "" && s.Http.KeyFile != "" {
		tlsConfig, err := tls.NewServerTLSConfig(s.Http.CertFile, s.Http.KeyFile)
		if err != nil {
			log.Fatal(err)
		}
		s.Http.TLSConfig = tlsConfig
	}
	if s.HTTP3.Enable && s.HTTP3.CertFile != "" && s.HTTP3.KeyFile != "" {
		tlsConfig, err := tls.NewServerTLSConfig(s.HTTP3.CertFile, s.HTTP3.KeyFile)
		if err != nil {
			log.Fatal(err)
		}
		s.HTTP3.TLSConfig = tlsConfig
	}
	if s.BaseContext == nil {
		s.BaseContext = context.Background()
	}
	s.HttpOption.AccessLog = DefaultAccessLog
}

// implement initialize
func (s *Server) BeforeInject() {
	*s = *NewServer()
}
func (s *Server) AfterInject() {
	s.Init()
}

func (s *Server) WithOptions(options ...Option) *Server {
	for _, option := range options {
		option(s)
	}
	return s
}
