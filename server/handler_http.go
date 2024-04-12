package server

import (
	"bytes"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/cherry/context/http_context"
	httpi "github.com/hopeio/cherry/utils/net/http"
	gini "github.com/hopeio/cherry/utils/net/http/gin"
	"github.com/hopeio/cherry/utils/net/http/grpc/gateway"
	"io"

	stringsi "github.com/hopeio/cherry/utils/strings"
	"net/http"
)

func (s *Server) httpHandler() http.HandlerFunc {
	conf := s.Config
	enablePrometheus := conf.Prometheus
	// 默认使用gin
	ginServer := conf.Gin.New()
	s.GinHandler(ginServer)
	gini.Debug(ginServer)
	if len(s.Config.StaticFs) > 0 {
		for _, fs := range s.Config.StaticFs {
			ginServer.Static(fs.Prefix, fs.Root)
		}
	}

	if s.GraphqlHandler != nil {
		graphqlServer := handler.NewDefaultServer(s.GraphqlHandler)
		ginServer.Handle(http.MethodPost, "/api/graphql", func(ctx *gin.Context) {
			graphqlServer.ServeHTTP(ctx.Writer, ctx.Request)
		})
	}
	var gatewayServer http.Handler
	if s.GatewayHandler != nil {
		gatewayServer = gateway.Gateway(s.GatewayHandler)
		/*	ginServer.NoRoute(func(ctx *gin.Context) {
			gatewayServer.ServeHTTP(
				(*httpi.ResponseRecorder)(unsafe.Pointer(uintptr(*(*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(ctx))+8))))),
				ctx.Request)
			ctx.Writer.WriteHeader(http.StatusOK)
		})*/
	}

	// http.Handle("/", ginServer)
	var excludes = []string{"/api/v1/upload", "/api/v1/multiUpload", "/api/ws/chat"}
	var includes = []string{"/api"}
	return func(w http.ResponseWriter, r *http.Request) {
		// 暂时解决方法，三个路由
		if h, p := http.DefaultServeMux.Handler(r); p != "" {
			h.ServeHTTP(w, r)
			return
		}
		if !stringsi.HasPrefixes(r.RequestURI, includes) || stringsi.HasPrefixes(r.RequestURI, excludes) {
			ginServer.ServeHTTP(w, r)
			return
		}

		var body []byte
		if r.Method != http.MethodGet {
			body, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewReader(body))
		}
		recorder := httpi.NewRecorder(w.Header())

		ginServer.ServeHTTP(recorder, r)
		if recorder.Code == http.StatusNotFound && gatewayServer != nil {
			recorder.Reset()
			gatewayServer.ServeHTTP(recorder, r)
		}

		// 提取 recorder 中记录的状态码，写入到 ResponseWriter 中
		w.WriteHeader(recorder.Code)
		if recorder.Body != nil {
			// 将 recorder 记录的 Response Body 写入到 ResponseWriter 中，客户端收到响应报文体
			w.Write(recorder.Body.Bytes())
		}
		ctxi := http_context.ContextFromContext(r.Context())
		defaultAccessLog(ctxi, r.RequestURI, r.Method,
			stringsi.BytesToString(body), stringsi.BytesToString(recorder.Body.Bytes()),
			recorder.Code)
		if enablePrometheus {
			defaultPrometheusRecord(ctxi, r.RequestURI, r.Method, recorder.Code)
		}
	}
}
