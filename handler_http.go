/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package cherry

import (
	"bytes"
	"github.com/hopeio/context/httpctx"
	httpi "github.com/hopeio/utils/net/http"
	"github.com/hopeio/utils/net/http/consts"
	gini "github.com/hopeio/utils/net/http/gin"
	"github.com/hopeio/utils/net/http/gin/apidoc"
	"io"
	"strings"

	stringsi "github.com/hopeio/utils/strings"
	"net/http"
)

func (s *Server) httpHandler() http.HandlerFunc {

	//enablePrometheus := conf.EnablePrometheus
	// 默认使用gin
	ginServer := s.Gin.New()
	// TODO: 不记录日志
	if s.EnableApiDoc {
		apidoc.OpenApi(ginServer, s.ApiDocUriPrefix, s.ApiDocDir)
	}
	s.GinHandler(ginServer)
	if s.EnableDebugApi {
		gini.Debug(ginServer)
	}

	if len(s.HttpOption.StaticFs) > 0 {
		for _, fs := range s.HttpOption.StaticFs {
			ginServer.Static(fs.Prefix, fs.Root)
		}
	}

	// http.Handle("/", ginServer)
	var excludes = s.HttpOption.ExcludeLogPrefixes
	var includes = s.HttpOption.IncludeLogPrefixes
	return func(w http.ResponseWriter, r *http.Request) {
		for _, middlewares := range s.HttpOption.Middlewares {
			middlewares(w, r)
		}
		// 暂时解决方法，三个路由
		if h, p := http.DefaultServeMux.Handler(r); p != "" {
			h.ServeHTTP(w, r)
			return
		}
		// 不记录日志
		if len(excludes) > 0 {
			if stringsi.HasPrefixes(r.RequestURI, excludes) && !stringsi.HasPrefixes(r.RequestURI, includes) {
				ginServer.ServeHTTP(w, r)
				return
			}
		}

		var body []byte
		if r.Method != http.MethodGet {
			body, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewReader(body))
		}
		recorder := httpi.NewRecorder(w.Header())

		ginServer.ServeHTTP(recorder, r)

		// 提取 recorder 中记录的状态码，写入到 ResponseWriter 中
		w.WriteHeader(recorder.Code)
		if recorder.Body != nil {
			// 将 recorder 记录的 Response Body 写入到 ResponseWriter 中，客户端收到响应报文体
			w.Write(recorder.Body.Bytes())
		}
		ctxi, _ := httpctx.FromContextValue(r.Context())
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
	}
}
