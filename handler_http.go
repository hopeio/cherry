/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package cherry

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/hopeio/gox/context/httpctx"
	"github.com/hopeio/gox/errors"
	"github.com/hopeio/gox/log"
	httpx "github.com/hopeio/gox/net/http"
	"github.com/hopeio/gox/net/http/apidoc"
	"github.com/hopeio/gox/net/http/debug"
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
		debug.Handle(s.DebugHandler.UriPrefix)
	}
}

func (s *Server) httpHandler() http.Handler {

	var handler http.Handler
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.StackLogger().Errorw(fmt.Sprintf("panic: %v", err))
				w.Header().Set(httpx.HeaderContentType, gatewayx.Codec.ContentType(nil))

				se := &response.CommonResp{Code: int32(errors.Internal), Msg: sysErrMsg}
				buf, err := gatewayx.Codec.Marshal(se)
				if err != nil {
					log.Error(err)
				}
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

		var body []byte
		if r.Body != nil {
			body, _ = io.ReadAll(r.Body)
			r.Body.Close()
			r.Body = io.NopCloser(bytes.NewReader(body))
		}
		recorder := httpx.NewRecorder(w.Header())

		s.GinServer.ServeHTTP(recorder, r)

		// 提取 recorder 中记录的状态码，写入到 ResponseWriter 中
		w.WriteHeader(recorder.Code)
		if recorder.Body != nil {
			// 将 recorder 记录的 Response Body 写入到 ResponseWriter 中，客户端收到响应报文体
			w.Write(recorder.Body.Bytes())
		}
		ctxi, _ := httpctx.FromContext(r.Context())
		if s.AccessLog.RecordFunc != nil {
			s.AccessLog.RecordFunc(ctxi, &AccessLogParam{
				r.Method, r.RequestURI,
				Body{
					ContentType: r.Header.Get(httpx.HeaderContentType),
					Data:        body,
				},
				Body{
					ContentType: w.Header().Get(httpx.HeaderContentType),
					Data:        recorder.Body.Bytes(),
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

					handlerBack.Respond(w, r)
				})*/
		handler = otelhttp.NewHandler(handler, "server", s.Telemetry.otelhttpOpts...)
	}
	return handler
}
