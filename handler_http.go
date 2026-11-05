/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package mix

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hopeio/gox/log"
	httpx "github.com/hopeio/gox/net/http"
	"github.com/hopeio/gox/net/http/openapi"
	stringsx "github.com/hopeio/gox/strings"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

const (
	HeaderErrorCode    = "Error-Code"
	HeaderErrorMessage = "Error-Message"
)

// InternalHandler registers the health check, the OpenAPI document and the
// debug endpoints on the private mux of the internal port.
func (s *Server) InternalHandler(mux *http.ServeMux) {
	// Standard probe endpoints for k8s and load balancers.
	healthz := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	mux.Handle("/healthz", healthz)
	mux.Handle("/readyz", healthz)
	if s.Openapi.Enabled {
		openapi.Openapi(mux, s.Openapi.PathPrefix, s.Openapi.Dir)
	}
	if s.Debug.Enabled {
		httpx.HandleDebug(mux, s.Debug.PathPrefix)
	}
}

func (s *Server) httpHandler() http.Handler {
	var handler http.Handler
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.StackLogger().Errorw(fmt.Sprintf("panic: %v", err))
				code := strconv.Itoa(int(Internal))
				w.Header().Set(HeaderErrorCode, code)
				se := &ErrResp{Code: Internal, Msg: sysErrMsg}
				buf, contentType, _ := DefaultMarshal(r.Context(), se)
				w.Header().Set(httpx.HeaderContentType, contentType)
				// A panic is a server-side error, so the client has to see a 500
				// (a no-op once the response headers are already written).
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(buf)
			}
		}()
		// Do not log it.
		if len(s.AccessLog.ExcludePrefixes) > 0 {
			if stringsx.HasPrefixes(r.URL.Path, s.AccessLog.ExcludePrefixes) &&
				!stringsx.HasPrefixes(r.URL.Path, s.AccessLog.IncludePrefixes) {
				s.HttpHandler.ServeHTTP(w, r)
				return
			}
		}
		md := GetMetadata(r.Context())
		if md == nil {
			md = &Metadata{Request: r, ResponseWriter: w, RequestAt: time.Now()}
		}
		md.TraceId = trace.SpanFromContext(r.Context()).SpanContext().TraceID().String()
		md.Logger = log.DefaultLogger().With(zap.String(log.FieldTraceId, md.TraceId))
		md.Baggage = baggage.FromContext(r.Context())
		// gRPC metadata uses lowercase keys while HTTP headers use the
		// CanonicalMIME form, so the case has to be converted.
		md.IncomingMD = make(metadata.MD, len(r.Header))
		for k, v := range r.Header {
			md.IncomingMD[strings.ToLower(k)] = v
		}
		// A WebTransport CONNECT needs the original ResponseWriter: it
		// hijacks the connection, and the access log recorder only
		// implements http.ResponseWriter, which loses both the hijack and
		// the http3.Settingser the upgrade looks for.
		//
		// The session is not access logged either: it is a long-lived
		// stream, not a request with a response body.
		if isWebTransportRequest(r) {
			s.HttpHandler.ServeHTTP(w, r)
			return
		}
		recorder := httpx.NewRecorder(w, r)
		r.Body = &recorder.RequestRecorder
		s.HttpHandler.ServeHTTP(&recorder.ResponseRecorder, r)

		if s.AccessLog.RecordFunc != nil {
			recorder.RequestRecorder.ContentType = r.Header.Get(httpx.HeaderContentType)
			recorder.ResponseRecorder.ContentType = recorder.Header().Get(httpx.HeaderContentType)
			s.AccessLog.RecordFunc(r.Context(), &AccessLogParam{
				r.Method, r.RequestURI,
				recorder,
				md,
			})
		}
		recorder.Reset()
	})
	handler = httpx.UseMiddleware(handler, s.Middlewares...)
	if s.Otel.Enabled {
		return otelhttp.NewHandler(handler, "http", append([]otelhttp.Option{
			otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
				return r.RequestURI
			}),
			// Mirror gRPC PublicEndpointFn: untrusted callers get a fresh root
			// (linked to the client span); trusted internal calls inherit.
			otelhttp.WithPublicEndpointFn(func(r *http.Request) bool {
				return !s.IsInternalHTTPRequest(r)
			}),
		}, s.Otel.OtelhttpOpts...)...)
	}
	return handler
}
