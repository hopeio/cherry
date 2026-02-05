/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package cherry

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/hopeio/gox/context/httpctx"
	"github.com/hopeio/gox/errors"
	"github.com/hopeio/gox/log"
	httpx "github.com/hopeio/gox/net/http"
	"github.com/hopeio/gox/net/http/apidoc"
	gatewayx "github.com/hopeio/gox/net/http/grpc/gateway"
	stringsx "github.com/hopeio/gox/strings"
	"github.com/hopeio/protobuf/response"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func (s *Server) InternalHandler() {
	if s.ApiDoc.Enabled {
		apidoc.ApiDoc(http.DefaultServeMux, s.ApiDoc.UriPrefix, s.ApiDoc.Dir)
	}
	if s.Telemetry.Enabled && s.Telemetry.Prometheus.Enabled {
		http.Handle(s.Telemetry.Prometheus.HttpUri, promhttp.Handler())
	}
	if s.DebugHandler.Enabled {
		httpx.HandleDebug(s.DebugHandler.UriPrefix)
	}
}

func (s *Server) httpHandler() http.Handler {
	var handler http.Handler
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.StackLogger().Errorw(fmt.Sprintf("panic: %v", err))
				code := strconv.Itoa(int(errors.Internal))
				w.Header().Set(httpx.HeaderErrorCode, code)
				se := &response.ErrResp{Code: int32(errors.Internal), Msg: sysErrMsg}
				buf, contentType := gatewayx.DefaultMarshal(r.Context(), se)
				w.Header().Set(httpx.HeaderContentType, contentType)
				w.Write(buf)
			}
		}()
		// 不记录日志
		if len(s.AccessLog.ExcludePrefixes) > 0 {
			if stringsx.HasPrefixes(r.RequestURI, s.AccessLog.ExcludePrefixes) &&
				!stringsx.HasPrefixes(r.RequestURI, s.AccessLog.IncludePrefixes) {
				s.GinServer.ServeHTTP(w, r)
				return
			}
		}

		recorder := httpx.NewRecorder(w, r)
		r.Body = &recorder.RequestRecorder
		s.GinServer.ServeHTTP(&recorder.ResponseRecorder, r)
		ctxi, _ := httpctx.FromContext(r.Context())
		if s.AccessLog.RecordFunc != nil {
			recorder.RequestRecorder.ContentType = r.Header.Get(httpx.HeaderContentType)
			recorder.ResponseRecorder.ContentType = recorder.Header().Get(httpx.HeaderContentType)
			s.AccessLog.RecordFunc(ctxi, &AccessLogParam{
				r.Method, r.RequestURI,
				recorder,
			})
		}
		recorder.Reset()
	})
	if s.Telemetry.Enabled {
		return otelhttp.NewHandler(handler, "server", s.Telemetry.otelhttpOpts...)
	}
	return handler
}
