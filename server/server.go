package server

import (
	"context"
	"fmt"
	"github.com/hopeio/cherry/context/http_context"
	httpi "github.com/hopeio/cherry/utils/net/http"
	"github.com/hopeio/cherry/utils/net/http/grpc/web"
	stringsi "github.com/hopeio/cherry/utils/strings"
	"github.com/quic-go/quic-go"
	"github.com/rs/cors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
	"net"
	"net/http"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/cherry/utils/log"
	"github.com/hopeio/cherry/utils/net/http/grpc/gateway"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.uber.org/zap"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

type CustomContext func(c context.Context, r *http.Request) context.Context
type ConvertContext func(r *http.Request) *http_context.Context

func (s *Server) Start() {
	if s.Config == nil {
		s.Config = defaultServerConfig()
	}
	baseCtx := context.Background()
	if s.Config.BaseContext != nil {
		baseCtx = s.Config.BaseContext()
		if baseCtx == nil {
			log.Fatal("BaseContext returned a nil context")
		}
	}
	// Handle SIGINT (CTRL+C) gracefully.
	sigCtx, stop := signal.NotifyContext(baseCtx, // kill -SIGINT XXXX 或 Ctrl+c
		syscall.SIGINT, // register that too, it should be ok
		// os.Kill等同于syscall.Kill
		syscall.SIGKILL, // register that too, it should be ok
		// kill -SIGTERM XXXX
		syscall.SIGTERM,
	)
	defer stop()

	if s.OnBeforeStart != nil {
		s.OnBeforeStart(sigCtx)
	}

	grpcServer := s.grpcHandler()
	httpHandler := s.httpHandler()

	// cors
	if s.Config.EnableCors {
		var corsServer *cors.Cors
		if s.Config.Cors == nil {
			corsServer = cors.AllowAll()
		} else {
			corsServer = cors.New(*s.Config.Cors)
		}
		httpHandler = corsServer.Handler(httpHandler).(http.HandlerFunc)
	}

	// grpc-web
	var wrappedGrpc *web.WrappedGrpcServer
	if s.Config.EnableGrpcWeb {
		wrappedGrpc = web.WrapServer(grpcServer, s.Config.GrpcWebOption...)
	}

	enableTrace := s.Config.Trace

	//systemTracing := serviceConfig.SystemTracing
	if enableTrace {
		grpc.EnableTracing = true
		// Set up OpenTelemetry.

		otelShutdown, err := setupOTelSDK(sigCtx, enableTrace)
		if err != nil {
			log.Fatal(err)
		}
		// Handle shutdown properly so nothing leaks.
		defer func() {
			err = otelShutdown(context.Background())
			if err != nil {
				log.Error(err)
			}
		}()

	}
	var handler http.Handler
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Errorw(fmt.Sprintf("panic: %v", err), zap.String(log.FieldStack, stringsi.BytesToString(debug.Stack())))
				w.Header().Set(httpi.HeaderContentType, httpi.ContentJSONHeaderValue)
				_, err := w.Write(httpi.ResponseSysErr)
				if err != nil {
					log.Error(err)
				}
			}
		}()

		ctx, span := http_context.ContextFromRequest(http_context.RequestCtx{Request: r, Response: w}, enableTrace)

		r = r.WithContext(ctx.ContextWrapper())

		contentType := r.Header.Get(httpi.HeaderContentType)
		if strings.HasPrefix(contentType, httpi.ContentGRPCHeaderValue) {
			if strings.HasPrefix(contentType[len(httpi.ContentGRPCHeaderValue):], "-web") && wrappedGrpc != nil {
				wrappedGrpc.ServeHTTP(w, r)
			} else if r.ProtoMajor == 2 && grpcServer != nil {
				grpcServer.ServeHTTP(w, r) // gRPC Server
			}
		} else {
			httpHandler.ServeHTTP(w, r)
		}

		if span != nil {
			span.End()
		}
	})

	if enableTrace {
		http.DefaultClient = otelhttp.DefaultClient
		handlerBack := handler

		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiCounter.Add(r.Context(), 1)
			attr := semconv.HTTPRouteKey.String(r.RequestURI)

			span := trace.SpanFromContext(r.Context())
			span.SetAttributes(attr)

			labeler, _ := otelhttp.LabelerFromContext(r.Context())
			labeler.Add(attr)

			handlerBack.ServeHTTP(w, r)
		})
		handler = otelhttp.NewHandler(handler, "server")

	}
	server := &s.Config.Http
	server.BaseContext = func(_ net.Listener) context.Context {
		return sigCtx
	}
	// 为了提供grpc服务,默认启用http2
	h2Handler := h2c.NewHandler(handler, &s.Config.Http2)
	server.Handler = h2Handler
	// 服务注册
	//initialize.GlobalConfig.Register()

	srvErr := make(chan error, 1)
	if s.Config.Http3 != nil && s.Config.Http3.TLSConfig != nil {
		s.Config.Http3.Handler = handler
		s.Config.Http3.ConnContext = func(ctx context.Context, c quic.Connection) context.Context {
			return sigCtx
		}
		go func() {
			log.Infof("http3 listening: %s", s.Config.Http3.Addr)
			srvErr <- s.Config.Http3.ListenAndServe()
		}()
	}
	go func() {
		log.Infof("listening: %s", s.Config.Http.Addr)
		srvErr <- server.ListenAndServe()
	}()
	if s.OnAfterStart != nil {
		s.OnAfterStart(sigCtx)
	}
	// Wait for interruption.
	select {
	case err := <-srvErr:
		// Error when starting HTTP server.
		log.Fatalf("failed to serve: %v", err)
	case <-sigCtx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
		log.Debug("stop server")
	}

	if s.OnBeforeStop != nil {
		s.OnBeforeStop(sigCtx)
	}
	//服务关闭
	if grpcServer != nil {
		grpcServer.GracefulStop()
	}
	if err := server.Shutdown(context.Background()); err != nil {
		log.Error(err)
	}

	if s.OnAfterStop != nil {
		s.OnAfterStop(sigCtx)
	}
}

type Server struct {
	Config *Config
	// 注册 grpc 服务
	GRPCHandler func(*grpc.Server)
	// 注册 grpc-gateway 服务
	GatewayHandler gateway.GatewayHandler
	// 注册 gin 服务
	GinHandler func(*gin.Engine)
	// 注册 graphql 服务
	GraphqlHandler graphql.ExecutableSchema

	// 各种钩子函数
	OnBeforeStart func(context.Context)
	OnAfterStart  func(context.Context)
	OnBeforeStop  func(context.Context)
	OnAfterStop   func(context.Context)
}

func NewServer(config *Config, ginhandler func(*gin.Engine), grpchandler func(*grpc.Server), gatewayhandler gateway.GatewayHandler, graphqlhandler graphql.ExecutableSchema) *Server {
	return &Server{
		Config:         config,
		GinHandler:     ginhandler,
		GRPCHandler:    grpchandler,
		GatewayHandler: gatewayhandler,
		GraphqlHandler: graphqlhandler,
	}
}

func Start(s *Server) {
	s.Start()
}
