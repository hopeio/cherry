/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package cherry

import (
	"context"
	"fmt"
	"github.com/hopeio/context/httpctx"
	httpi "github.com/hopeio/utils/net/http"
	"github.com/hopeio/utils/net/http/grpc/web"
	stringsi "github.com/hopeio/utils/strings"
	"github.com/quic-go/quic-go"
	"github.com/rs/cors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net"
	"net/http"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/hopeio/utils/log"
	"go.uber.org/zap"
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

	if s.OnStart != nil {
		s.OnStart(sigCtx)
	}

	grpcServer := s.grpcHandler()
	httpHandler := s.httpHandler()

	// cors
	if s.EnableCors {
		var corsServer *cors.Cors
		if s.Cors == nil {
			corsServer = cors.AllowAll()
		} else {
			corsServer = cors.New(*s.Cors)
		}
		httpHandler = corsServer.Handler(httpHandler).(http.HandlerFunc)
	}

	// grpc-web
	var wrappedGrpc *web.WrappedGrpcServer
	if s.EnableGrpcWeb {
		wrappedGrpc = web.WrapServer(grpcServer, s.GrpcWebOptions...)
	}

	enableTelemetry := s.EnableTelemetry

	//systemTracing := serviceConfig.SystemTracing
	if enableTelemetry {
		grpc.EnableTracing = true
		http.DefaultClient = otelhttp.DefaultClient
		// Set up OpenTelemetry.

		otelShutdown, err := setupOTelSDK(sigCtx, &s.TelemetryConfig)
		if err != nil {
			log.Fatal(err)
		}
		// Handle shutdown properly so nothing leaks.
		defer otelShutdown(sigCtx)

	}

	var handler http.Handler
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Errorw(fmt.Sprintf("panic: %v", err), zap.String(log.FieldStack, stringsi.BytesToString(debug.Stack())))
				w.Header().Set(httpi.HeaderContentType, httpi.ContentTypeJson)
				_, err := w.Write(httpi.ResponseSysErr)
				if err != nil {
					log.Error(err)
				}
			}
		}()

		// 简单的中间件支持
		for _, middleware := range s.Middlewares {
			middleware(w, r)
		}

		ctx := httpctx.FromRequest(httpctx.RequestCtx{Request: r, Response: w})

		r = r.WithContext(ctx.Wrapper())

		contentType := r.Header.Get(httpi.HeaderContentType)
		if strings.HasPrefix(contentType, httpi.ContentTypeGrpc) {
			if strings.HasPrefix(contentType[len(httpi.ContentTypeGrpc):], "-web") && wrappedGrpc != nil {
				wrappedGrpc.ServeHTTP(w, r)
			} else if r.ProtoMajor == 2 && grpcServer != nil {
				grpcServer.ServeHTTP(w, r) // gRPC Server
			}
		} else {
			httpHandler.ServeHTTP(w, r)
		}

		ctx.RootSpan().End()
	})

	if enableTelemetry {

		/*		handlerBack := handler

				handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					//apiCounter.Add(r.Context(), 1)
					attr := semconv.HTTPRouteKey.String(r.RequestURI)

					span := trace.SpanFromContext(r.Context())
					span.SetAttributes(attr)

					labeler, _ := otelhttp.LabelerFromContext(r.Context())
					labeler.Add(attr)

					handlerBack.ServeHTTP(w, r)
				})*/
		handler = otelhttp.NewHandler(handler, "server")
	}
	server := &s.Http
	server.BaseContext = func(_ net.Listener) context.Context {
		return sigCtx
	}
	// 为了提供grpc服务,默认启用http2
	h2Handler := h2c.NewHandler(handler, &s.Http2)
	server.Handler = h2Handler

	srvErr := make(chan error, 1)
	if s.Http3 != nil && s.Http3.TLSConfig != nil {
		s.Http3.Handler = handler
		s.Http3.ConnContext = func(ctx context.Context, c quic.Connection) context.Context {
			return sigCtx
		}
		go func() {
			log.Infof("http3 listening: %s", s.Http3.Addr)
			srvErr <- s.Http3.ListenAndServe()
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
	if err := server.Shutdown(context.Background()); err != nil {
		log.Error(err)
	}

	if s.OnStop != nil {
		s.OnStop(sigCtx)
	}
}
