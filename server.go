/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package mix

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptrace"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/hopeio/gox/log"
	httpx "github.com/hopeio/gox/net/http"
	"github.com/quic-go/quic-go"
	"github.com/rs/cors"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
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
	// BaseContext is the root context (settable through WithContext); it is
	// cancelled together with the signal context.
	baseCtx := s.BaseContext
	if baseCtx == nil {
		baseCtx = context.Background()
	}
	// Handler SIGINT (CTRL+C) gracefully.
	// Note: SIGKILL cannot be caught, so there is no need to register it.
	sigCtx, stop := signal.NotifyContext(baseCtx, // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGINT,
		// kill -SIGTERM XXXX
		syscall.SIGTERM,
	)
	defer stop()

	grpcServer := s.grpcHandler()
	httpHandler := s.httpHandler()

	// cors
	if s.Cors.Enabled {
		httpHandler = cors.New(s.Cors.Options).Handler(httpHandler)
	}

	// Set up OpenTelemetry instrumentation (HTTP/gRPC wrappers).
	var otelShutdown func(context.Context) error
	if s.Otel.Enabled {
		http.DefaultClient = &http.Client{
			Transport: otelhttp.NewTransport(
				http.DefaultTransport,
				otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
					return otelhttptrace.NewClientTrace(ctx)
				}),
			),
		}
		shutdownFunc, err := setupOTelSDK(sigCtx, &s.Otel)
		if err != nil {
			log.Fatal(err)
		}
		if shutdownFunc != nil {
			otelShutdown = shutdownFunc
			// sigCtx is already cancelled by the time this defer runs. An
			// exporter given a cancelled context drops the final flush, so
			// use a fresh bounded context while keeping sigCtx's values.
			defer func() {
				flushCtx, cancel := context.WithTimeout(context.WithoutCancel(sigCtx), 5*time.Second)
				defer cancel()
				if err := shutdownFunc(flushCtx); err != nil {
					log.Error(err)
				}
			}()
		}
		s.tracer = otel.Tracer(ScopeName)
		s.meter = otel.Meter(ScopeName)
	}

	onePort := s.Grpc.Addr == "" || s.Grpc.Addr == s.Addr
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		md := Metadata{
			Request:        r,
			ResponseWriter: w,
			RequestAt:      time.Now(),
		}
		r = r.WithContext(WithMetadata(r.Context(), &md))
		if onePort && strings.HasPrefix(r.Header.Get(httpx.HeaderContentType), httpx.ContentTypeGrpc) {
			if r.ProtoMajor == 2 && grpcServer != nil {
				md.RequestType = RequestTypeGrpc
				grpcServer.ServeHTTP(w, r)
			} else {
				http.NotFound(w, r)
			}
		} else {
			httpHandler.ServeHTTP(w, r)
		}
	})

	s.Server.Handler = handler

	if s.Server.BaseContext == nil {
		s.Server.BaseContext = func(_ net.Listener) context.Context {
			return sigCtx
		}
	}

	// HTTP/2 is enabled by default because the gRPC service needs it.
	if s.Server.TLSConfig == nil {
		s.Server.Protocols = new(http.Protocols)
		s.Server.Protocols.SetHTTP1(true)
		s.Server.Protocols.SetUnencryptedHTTP2(true)
	}

	// The four listener goroutines (grpc / http3 / http / internal) can all
	// report an error.
	srvErr := make(chan error, 4)

	if s.Grpc.Addr != "" && s.Grpc.Addr != s.Addr {
		go func() {
			log.Infof("grpc listening: %s", s.Grpc.Addr)
			listener, err := net.Listen("tcp", s.Grpc.Addr)
			if err != nil {
				srvErr <- err
				return
			}
			srvErr <- grpcServer.Serve(listener)
		}()
	}

	// HTTP/3 is an optional extra listener over the same handler, so it gets its own
	// preparation pass before the block below decides whether to run it.
	s.prepareHTTP3()

	if s.HTTP3.Enabled {
		s.HTTP3.Handler = handler
		if s.HTTP3.ConnContext == nil {
			s.HTTP3.ConnContext = func(connCtx context.Context, _ *quic.Conn) context.Context {
				return bridgeShutdownContext(sigCtx, connCtx)
			}
		}
		go func() {
			log.Infof("http3 listening: %s", s.HTTP3.Addr)
			srvErr <- s.serveHTTP3(s.HTTP3.CertFile, s.HTTP3.KeyFile, handler)
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

	if !s.DisableInternalServer {
		s.InternalHandler(http.DefaultServeMux)
		if s.InternalServer.BaseContext == nil {
			s.InternalServer.BaseContext = func(_ net.Listener) context.Context {
				return sigCtx
			}
		}
		if s.InternalServer.Handler == nil {
			s.InternalServer.Handler = http.DefaultServeMux
		}
		go func() {
			log.Infof("internal listening: %s", s.InternalServer.Addr)
			srvErr <- s.InternalServer.ListenAndServe()
		}()
	}

	// Wait for interruption.
	select {
	case err := <-srvErr:
		// Startup or listen failure: log.Fatal skips every defer, so do the
		// cleanup by hand first (signal deregistration, OTel flush).
		stop()
		if otelShutdown != nil {
			// Match the deferred path: baseCtx may carry values supplied
			// through WithContext, and the final flush should keep them.
			flushCtx, cancel := context.WithTimeout(context.WithoutCancel(sigCtx), 5*time.Second)
			if err := otelShutdown(flushCtx); err != nil {
				log.Error(err)
			}
			cancel()
		}
		log.Fatalf("failed to serve: %v", err)
	case <-sigCtx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
		log.Debug("stop server")
	}

	// Shutting down: sigCtx is already cancelled and cannot drive the
	// graceful shutdown, so build a new context with a timeout.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// gRPC / HTTP / HTTP3 / internal shut down gracefully in parallel,
	// sharing one timeout budget.
	var wg sync.WaitGroup
	if grpcServer != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			grpcStopped := make(chan struct{})
			go func() {
				grpcServer.GracefulStop()
				close(grpcStopped)
			}()
			select {
			case <-grpcStopped:
			case <-shutdownCtx.Done():
				// Force the shutdown once the timeout expires.
				grpcServer.Stop()
			}
		}()
	}
	if s.HTTP3.Enabled {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.HTTP3.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Error(err)
			}
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error(err)
		}
	}()
	if !s.DisableInternalServer {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.InternalServer.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Error(err)
			}
		}()
	}
	wg.Wait()
}

func (s *Server) WithContext(ctx context.Context) *Server {
	s.BaseContext = ctx
	return s
}

// bridgeShutdownContext derives the context for one HTTP/3 connection.
//
// quic-go builds connCtx from conn.Context() and attaches ServerContextKey,
// LocalAddrContextKey and RemoteAddrContextKey before calling ConnContext.
// It also cancels that context when the QUIC connection ends, so requests
// still in flight on the connection stop with it. Returning the process
// shutdown context instead of deriving from connCtx would throw all of that
// away: those values disappear, and a client disconnect no longer cancels
// its requests.
//
// A context has one parent, so the shutdown signal is bridged in with
// AfterFunc. Both cancel paths are covered: the QUIC connection ending
// cancels connCtx (and therefore the derived ctx), and a SIGTERM cancels
// shutdown, which cancels the derived ctx through the callback. The second
// AfterFunc releases the registration as soon as either one fires, so a
// long-lived process does not accumulate one callback per connection.
func bridgeShutdownContext(shutdown, connCtx context.Context) context.Context {
	ctx, cancel := context.WithCancel(connCtx)
	stop := context.AfterFunc(shutdown, cancel)
	context.AfterFunc(ctx, func() {
		stop()
		cancel()
	})
	return ctx
}

// prepareHTTP3 fills the HTTP/3 defaults that belong to this library instead of being
// re-implemented by every consumer, and turns "no TLS material" into a disabled option.
//
// The certificate pair falls back to the main server's: HTTP/3 and HTTPS are the same
// endpoint over a different transport, and demanding both be configured only produced
// deployments where h3 — and WebTransport, which rides on this same QUIC stack — stayed
// silently off while the API answered fine.
//
// A missing pair must not fail the process: srvErr is consumed by log.Fatalf, so an
// optional transport would otherwise outrank the listener that is already serving. It
// disables itself with an error line instead, and only ever disables — never enables.
//
// Addr mirrors what quic-go binds when it is empty (":https", the same port), so the
// startup line names the port actually open.
func (s *Server) prepareHTTP3() {
	if !s.HTTP3.Enabled {
		return
	}
	if s.HTTP3.CertFile == "" {
		s.HTTP3.CertFile = s.CertFile
	}
	if s.HTTP3.KeyFile == "" {
		s.HTTP3.KeyFile = s.KeyFile
	}
	if s.HTTP3.Addr == "" {
		s.HTTP3.Addr = ":443"
	}
	if s.HTTP3.CertFile == "" || s.HTTP3.KeyFile == "" {
		s.HTTP3.Enabled = false
		log.Errorf("mix: http3 disabled (WebTransport with it): QUIC requires TLS, set HTTP3.CertFile/HTTP3.KeyFile or the server's CertFile/KeyFile")
	}
}
