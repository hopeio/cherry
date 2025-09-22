/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package cherry

import (
	"bytes"
	"github.com/hopeio/context/httpctx"
	httpx "github.com/hopeio/gox/net/http"
	"github.com/hopeio/gox/net/http/apidoc"
	"github.com/hopeio/gox/net/http/consts"
	"github.com/hopeio/gox/net/http/debug"
	stringsx "github.com/hopeio/gox/strings"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"io"
	"net/http"
	"strings"
)

func (s *Server) InternalHandler() {
	if s.ApiDoc.Enabled {
		apidoc.OpenApi(http.DefaultServeMux, s.ApiDoc.UriPrefix, s.ApiDoc.Dir)
	}
	if s.Telemetry.Enabled && s.Telemetry.EnablePrometheus {
		http.Handle(s.Telemetry.PromHttpUri, promhttp.Handler())
	}
	if s.DebugHandler.Enabled {
		debug.Handle(s.DebugHandler.UriPrefix)
	}
}

func (s *Server) httpHandler() http.Handler {
	s.InternalHandler()
	ginServer := s.Gin.New()
	s.GinHandler(ginServer)

	var excludes = s.HttpOption.ExcludeLogPrefixes
	var includes = s.HttpOption.IncludeLogPrefixes
	var handler http.Handler
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 不记录日志
		if len(excludes) > 0 {
			if stringsx.HasPrefixes(r.RequestURI, excludes) && !stringsx.HasPrefixes(r.RequestURI, includes) {
				ginServer.ServeHTTP(w, r)
				return
			}
		}

		var body []byte
		if r.Method != http.MethodGet {
			body, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewReader(body))
		}
		recorder := httpx.NewRecorder(w.Header())

		ginServer.ServeHTTP(recorder, r)

		// 提取 recorder 中记录的状态码，写入到 ResponseWriter 中
		w.WriteHeader(recorder.Code)
		if recorder.Body != nil {
			// 将 recorder 记录的 Response Body 写入到 ResponseWriter 中，客户端收到响应报文体
			w.Write(recorder.Body.Bytes())
		}
		ctxi, _ := httpctx.FromContext(r.Context())
		if s.HttpOption.AccessLog != nil {
			s.HttpOption.AccessLog(ctxi, &AccessLogParam{
				r.Method, r.RequestURI,
				Body{
					IsJson: strings.HasPrefix(r.Header.Get(consts.HeaderContentType), consts.ContentTypeJson),
					Data:   body,
				},
				Body{
					IsJson: strings.HasPrefix(w.Header().Get(consts.HeaderContentType), consts.ContentTypeJson),
					Data:   recorder.Body.Bytes(),
				},
				recorder.Code,
			})
		}
		/*		if enablePrometheus {
				defaultMetricsRecord(ctxi, r.RequestURI, r.Method, recorder.Code)
			}*/
	})
	if s.Telemetry.Enabled {

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
		handler = otelhttp.NewHandler(handler, "server", s.Telemetry.otelhttpOpts...)
	}
	return handler
}
