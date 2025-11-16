/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package cherry

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"reflect"
	"strings"
	"syscall"

	"github.com/hopeio/context/httpctx"
	"github.com/hopeio/gox/crypto/tls"
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

func (s *Server) Run() {
	s.Init()
	baseCtx := context.Background()

	// Handle SIGINT (CTRL+C) gracefully.
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
		// Handle shutdown properly so nothing leaks.
		defer otelShutdown(sigCtx)
	}

	var handler http.Handler
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h, p := http.DefaultServeMux.Handler(r); p != "" {
			h.ServeHTTP(w, r)
			return
		}
		defer func() {
			if err := recover(); err != nil {
				log.StackLogger().Errorw(fmt.Sprintf("panic: %v", err))
				w.Header().Set(httpx.HeaderContentType, httpx.ContentTypeJson)
				_, err := w.Write(httpx.RespSysErr)
				if err != nil {
					log.Error(err)
				}
			}
		}()

		// 简单的中间件支持
		for _, middleware := range s.Middlewares {
			middleware(w, r)
		}

		ctx := httpctx.FromRequest(httpctx.RequestCtx{Request: r, ResponseWriter: w})

		r = r.WithContext(ctx.Wrapper())

		contentType := r.Header.Get(httpx.HeaderContentType)
		if strings.HasPrefix(contentType, httpx.ContentTypeGrpc) {
			if strings.HasPrefix(contentType[len(httpx.ContentTypeGrpc):], "-web") && wrappedGrpc != nil {
				wrappedGrpc.ServeHTTP(w, r)
			} else if r.ProtoMajor == 2 && grpcServer != nil {
				grpcServer.ServeHTTP(w, r) // gRPC Server
			}
		} else {
			httpHandler.ServeHTTP(w, r)
		}

		ctx.RootSpan().End()
	})

	server := &s.Http
	server.BaseContext = func(_ net.Listener) context.Context {
		return sigCtx
	}
	// 为了提供grpc服务,默认启用http2
	if s.TLSConfig != nil {
		err := http2.ConfigureServer(&server.Server, &s.HTTP2)
		if err != nil {
			log.Fatal(err)
		}
		server.Handler = handler
	} else {
		h2Handler := h2c.NewHandler(handler, &s.HTTP2)
		server.Handler = h2Handler
	}
	srvErr := make(chan error, 1)
	if s.HTTP3.Enabled {
		if s.HTTP3.TLSConfig == nil {
			if s.HTTP3.CertFile != "" && s.HTTP3.KeyFile != "" {
				var err error
				s.HTTP3.TLSConfig, err = tls.NewServerTLSConfig(s.HTTP3.CertFile, s.HTTP3.KeyFile)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				log.Fatal("http3 need certFile and keyFile")
			}
		}

		s.HTTP3.Handler = handler
		s.HTTP3.ConnContext = func(ctx context.Context, c quic.Connection) context.Context {
			return sigCtx
		}
		go func() {
			log.Infof("http3 listening: %s", s.HTTP3.Addr)
			srvErr <- s.HTTP3.ListenAndServe()
		}()
	}
	go func() {
		log.Infof("listening: %s", s.Http.Addr)
		srvErr <- server.ListenAndServe()
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
	if err := server.Shutdown(sigCtx); err != nil {
		log.Error(err)
	}
}
