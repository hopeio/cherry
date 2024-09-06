package cherry

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/utils/crypto/tls"
	"github.com/hopeio/utils/log"
	gini "github.com/hopeio/utils/net/http/gin"
	"github.com/hopeio/utils/net/http/grpc/gateway/grpc-gateway"
	"github.com/hopeio/utils/net/http/grpc/web"
	"github.com/hopeio/utils/validation/validator"
	"github.com/quic-go/quic-go/http3"
	"github.com/rs/cors"
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
	c.StopTimeout = 5 * time.Second
	gin.SetMode(gin.ReleaseMode)
	gin.DisableBindValidation()
	validator.DefaultValidator = nil // 自己做校验
	c.EnableCors = true
	c.EnableTelemetry = true
	c.MetricsInterval = time.Minute
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
}
type Http3 struct {
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
	Name        string
	Http        Http
	Http2       http2.Server
	Http3       *Http3
	StopTimeout time.Duration
	Gin         gini.Config
	EnableCors  bool
	Cors        *cors.Options
	Middlewares []http.HandlerFunc
	HttpOption  HttpOption
	// Grpc options
	GrpcOptions                                   []grpc.ServerOption
	EnableGrpcWeb                                 bool
	GrpcWebOptions                                []web.Option
	EnableTelemetry, EnableDebugApi, EnableApiDoc bool
	ApiDocUriPrefix, ApiDocDir                    string
	TelemetryConfig
	BaseContext context.Context
	// 注册 grpc 服务
	GrpcHandler func(*grpc.Server)
	// 注册 grpc-gateway 服务
	GatewayHandler grpc_gateway.GatewayHandler
	// 注册 gin 服务
	GinHandler func(*gin.Engine)
	// 注册 graphql 服务
	GraphqlHandler graphql.ExecutableSchema
	// 各种钩子函数
	OnStart func(context.Context)
	OnStop  func(context.Context)
}

type TelemetryConfig struct {
	EnablePrometheus bool
	MetricsInterval  time.Duration
	propagator       propagation.TextMapPropagator
	tracerProvider   *sdktrace.TracerProvider
	meterProvider    *sdkmetric.MeterProvider
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
	if s.Http3 != nil && s.Http3.Addr == "" {
		s.Http3.Addr = ":8080"
	}
	log.DurationNotify("ReadTimeout", s.Http.ReadTimeout, time.Second)
	log.DurationNotify("WriteTimeout", s.Http.WriteTimeout, time.Second)
	if s.StopTimeout == 0 {
		s.StopTimeout = 5 * time.Second
	}
	log.DurationNotify("StopTimeout", s.StopTimeout, time.Second)
	if s.Http.CertFile != "" && s.Http.KeyFile != "" {
		tlsConfig, err := tls.NewServerTLSConfig(s.Http.CertFile, s.Http.KeyFile)
		if err != nil {
			log.Fatal(err)
		}
		s.Http.TLSConfig = tlsConfig
	}
	if s.Http3 != nil && s.Http3.CertFile != "" && s.Http3.KeyFile != "" {
		tlsConfig, err := tls.NewServerTLSConfig(s.Http3.CertFile, s.Http3.KeyFile)
		if err != nil {
			log.Fatal(err)
		}
		s.Http3.TLSConfig = tlsConfig
	}
	if s.BaseContext == nil {
		s.BaseContext = context.Background()
	}
	s.HttpOption.AccessLog = defaultAccessLog
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
