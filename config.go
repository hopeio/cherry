/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package cherry

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hopeio/gox/crypto/tls"
	"github.com/hopeio/gox/log"
	httpx "github.com/hopeio/gox/net/http"
	"github.com/hopeio/gox/net/http/grpc/web"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/quic-go/quic-go/http3"
	"github.com/rs/cors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
)

type Http3Config struct {
	Enabled bool
	http3.Server
	CertFile string
	KeyFile  string
}

type AccessLogConfig struct {
	RecordFunc      AccessLog
	ExcludePrefixes []string
	IncludePrefixes []string
}

type Server struct {
	http.Server
	CertFile       string
	KeyFile        string
	AccessLog      AccessLogConfig
	HTTP2          Http2Config
	HTTP3          Http3Config
	Cors           CorsConfig
	Grpc           GrpcConfig
	InternalServer http.Server
	ApiDoc         ApiDocConfig
	Otel           OtelConfig
	Prometheus     PrometheusConfig
	DebugHandler   DebugHandlerConfig
	BaseContext    context.Context
	Middlewares    []httpx.Middleware
	GinServer      *gin.Engine
	GrpcHandler    func(*grpc.Server)
}

type Http2Config struct {
	NewWriteScheduler func() http2.WriteScheduler
}

type DebugHandlerConfig struct {
	Enabled   bool
	UriPrefix string
}

type ApiDocConfig struct {
	Enabled        bool
	UriPrefix, Dir string
}

type GrpcConfig struct {
	RecordFunc               GrpcAccessLog
	EnableGrpcWeb            bool
	GrpcWebOptions           []web.Option
	Options                  []grpc.ServerOption
	UnaryServerInterceptors  []grpc.UnaryServerInterceptor
	StreamServerInterceptors []grpc.StreamServerInterceptor
}

type CorsConfig struct {
	Enabled bool
	cors.Options
}

type OtelConfig struct {
	Enabled      bool
	OtelhttpOpts []otelhttp.Option
	OtelgrpcOpts []otelgrpc.Option
}

func (c *OtelConfig) SetOtelhttpHandlerOpts(otelhttpOpts []otelhttp.Option) {
	c.OtelhttpOpts = otelhttpOpts
}

func (c *OtelConfig) SetOtelgrpcOptsHandlerOpts(otelgrpcOpts []otelgrpc.Option) {
	c.OtelgrpcOpts = otelgrpcOpts
}

type PrometheusConfig struct {
	Enabled bool
	HttpURI string
	promhttp.HandlerOpts
}

func (s *Server) Init() {
	gin.SetMode(gin.ReleaseMode)
	if s.BaseContext == nil {
		s.BaseContext = context.Background()
	}
	if s.Addr == "" {
		s.Addr = ":8080"
	}

	if s.AccessLog.RecordFunc == nil {
		s.AccessLog.RecordFunc = DefaultAccessLog
	}
	if s.Grpc.RecordFunc == nil {
		s.Grpc.RecordFunc = DefaultGrpcAccessLog
	}

	if s.GinServer == nil {
		s.GinServer = gin.New()
	}

	if s.InternalServer.Addr == "" {
		s.InternalServer.Addr = ":8081"
	}

	log.ValueLevelNotify("ReadTimeout", s.ReadTimeout, time.Second)
	log.ValueLevelNotify("WriteTimeout", s.WriteTimeout, time.Second)
	if s.CertFile != "" && s.KeyFile != "" {
		tlsConfig, err := tls.NewServerTLSConfig(s.CertFile, s.KeyFile)
		if err != nil {
			log.Fatal(err)
		}
		s.TLSConfig = tlsConfig
	}
	if s.HTTP3.Enabled {
		if s.HTTP3.Addr == "" {
			s.HTTP3.Addr = ":8080"
		}
		if s.HTTP3.CertFile != "" && s.HTTP3.KeyFile != "" {
			tlsConfig, err := tls.NewServerTLSConfig(s.HTTP3.CertFile, s.HTTP3.KeyFile)
			if err != nil {
				log.Fatal(err)
			}
			s.HTTP3.TLSConfig = tlsConfig
		}
	}
	if s.Cors.Enabled {
		if len(s.Cors.AllowedOrigins) == 0 {
			s.Cors.AllowedOrigins = []string{"*"}
		}
		if len(s.Cors.AllowedMethods) == 0 {
			s.Cors.AllowedMethods = []string{http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete}
		}
		if len(s.Cors.AllowedHeaders) == 0 {
			s.Cors.AllowedHeaders = []string{"*"}
		}
	}

	if s.Prometheus.Enabled {
		if s.Prometheus.HttpURI == "" {
			s.Prometheus.HttpURI = "/metrics"
		}
	}

}

// implement initialize
func (s *Server) BeforeInject() {
	s.Init()
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
