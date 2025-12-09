/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package cherry

import (
	"context"
	"net"
	"net/http"
	"os/signal"
	"reflect"
	"strings"
	"syscall"

	"github.com/hopeio/gox/context/httpctx"
	"github.com/hopeio/gox/log"
	httpx "github.com/hopeio/gox/net/http"
	"github.com/hopeio/gox/net/http/grpc/web"
	"github.com/quic-go/quic-go"
	"github.com/rs/cors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

func NewServer(options ...Option) *Server {
	s := &Server{}
	s.Init()
	for _, option := range options {
		option(s)
	}
	return s
}

func (s *Server) Run() {
	baseCtx := context.Background()
	// Handler SIGINT (CTRL+C) gracefully.
	sigCtx, stop := signal.NotifyContext(baseCtx, // kill -SIGINT XXXX 或 Ctrl+c
		syscall.SIGINT, // register that too, it should be ok
		// os.Kill等同于syscall.Kill
		syscall.SIGKILL, // register that too, it should be ok
		// kill -SIGTERM XXXX
		syscall.SIGTERM,
	)
	defer stop()

	grpcServer := s.grpcHandler()
	httpHandler := s.httpHandler()

	// cors
	if s.Cors.Enabled {
		var corsServer *cors.Cors
		if reflect.ValueOf(&s.Cors.Options).Elem().IsZero() {
			corsServer = cors.AllowAll()
		} else {
			corsServer = cors.New(s.Cors.Options)
		}
		httpHandler = corsServer.Handler(httpHandler)
	}

	// grpc-web
	var wrappedGrpc *web.WrappedGrpcServer
	if s.Grpc.EnableGrpcWeb {
		wrappedGrpc = web.WrapServer(grpcServer, s.Grpc.GrpcWebOptions...)
	}

	// Set up OpenTelemetry.
	if s.Telemetry.Enabled {

		grpc.EnableTracing = true
		http.DefaultClient = otelhttp.DefaultClient

		otelShutdown, err := s.Telemetry.setupOTelSDK(sigCtx)
		if err != nil {
			log.Fatal(err)
		}
		// Handler shutdown properly so nothing leaks.
		defer otelShutdown(sigCtx)
	}

	mwHandler := httpx.UseMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get(httpx.HeaderContentType)
		if strings.HasPrefix(contentType, httpx.ContentTypeGrpc) {
			if strings.HasPrefix(contentType[len(httpx.ContentTypeGrpc):], "-web") && wrappedGrpc != nil {
				wrappedGrpc.ServeHTTP(w, r)
			} else if r.ProtoMajor == 2 && grpcServer != nil {
				grpcServer.ServeHTTP(w, r) // gRPC Server
			} else {
				http.NotFound(w, r)
			}
		} else {
			httpHandler.ServeHTTP(w, r)
		}
	}), s.Middlewares...)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := httpctx.FromRequest(httpctx.RequestCtx{Request: r, ResponseWriter: w})
		r = r.WithContext(ctx.Wrapper())
		mwHandler.ServeHTTP(w, r)
		ctx.RootSpan().End()
	})

	if s.Server.BaseContext == nil {
		s.Server.BaseContext = func(_ net.Listener) context.Context {
			return sigCtx
		}
	}

	// 为了提供grpc服务,默认启用http2
	if s.Server.TLSConfig != nil || (s.CertFile != "" && s.KeyFile != "") {
		err := http2.ConfigureServer(&s.Server, &s.HTTP2)
		if err != nil {
			log.Fatal(err)
		}
		s.Server.Handler = handler
	} else {
		h2Handler := h2c.NewHandler(handler, &s.HTTP2)
		s.Server.Handler = h2Handler
	}
	srvErr := make(chan error, 1)
	if s.HTTP3.Enabled {
		s.HTTP3.Handler = handler
		if s.HTTP3.ConnContext == nil {
			s.HTTP3.ConnContext = func(ctx context.Context, c *quic.Conn) context.Context {
				return sigCtx
			}
		}
		go func() {
			log.Infof("http3 listening: %s", s.HTTP3.Addr)
			if s.HTTP3.CertFile != "" && s.HTTP3.KeyFile != "" {
				srvErr <- s.HTTP3.ListenAndServeTLS(s.CertFile, s.KeyFile)
			} else {
				srvErr <- s.HTTP3.ListenAndServe()
			}
		}()
	}
	go func() {
		log.Infof("listening: %s", s.Addr)
		if s.CertFile != "" && s.KeyFile != "" {
			srvErr <- s.ListenAndServeTLS(s.CertFile, s.KeyFile)
		} else {
			srvErr <- s.ListenAndServe()
		}
	}()

	go func() {
		log.Infof("internal listening: %s", s.InternalServer.Addr)
		s.InternalHandler()
		if s.InternalServer.BaseContext == nil {
			s.InternalServer.BaseContext = func(_ net.Listener) context.Context {
				return sigCtx
			}
			s.InternalServer.Handler = http.DefaultServeMux
		}
		srvErr <- s.InternalServer.ListenAndServe()
	}()

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

	//服务关闭
	if grpcServer != nil {
		grpcServer.GracefulStop()
	}
	if err := s.Shutdown(sigCtx); err != nil {
		log.Error(err)
	}
}
