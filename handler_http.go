package cherry

import (
	"bytes"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hopeio/context/httpctx"
	httpi "github.com/hopeio/utils/net/http"
	gini "github.com/hopeio/utils/net/http/gin"
	"github.com/hopeio/utils/net/http/grpc/gateway/grpc-gateway"
	"io"

	stringsi "github.com/hopeio/utils/strings"
	"net/http"
)

func (s *Server) httpHandler() http.HandlerFunc {

	//enablePrometheus := conf.EnablePrometheus
	// 默认使用gin
	ginServer := s.Gin.New()
	// TODO: 不记录日志
	if s.EnableApiDoc {
		gini.OpenApi(ginServer, s.ApiDocUriPrefix, s.ApiDocDir)
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

	if s.GraphqlHandler != nil {
		graphqlServer := handler.NewDefaultServer(s.GraphqlHandler)
		ginServer.Handle(http.MethodPost, "/api/graphql", func(ctx *gin.Context) {
			graphqlServer.ServeHTTP(ctx.Writer, ctx.Request)
		})
	}
	var gatewayServer *runtime.ServeMux
	if s.GatewayHandler != nil {
		gatewayServer = grpc_gateway.New()
		s.GatewayHandler(s.BaseContext, gatewayServer)
		/*	ginServer.NoRoute(func(ctx *gin.Context) {
			gatewayServer.ServeHTTP(
				(*httpi.ResponseRecorder)(unsafe.Pointer(uintptr(*(*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(ctx))+8))))),
				ctx.Request)
			ctx.Writer.WriteHeader(http.StatusOK)
		})*/
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
		if len(includes) > 0 || len(excludes) > 0 {
			if !stringsi.HasPrefixes(r.RequestURI, includes) || stringsi.HasPrefixes(r.RequestURI, excludes) {
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
		ctxi := httpctx.FromContextValue(r.Context())
		defaultAccessLog(ctxi, r.RequestURI, r.Method,
			stringsi.BytesToString(body), stringsi.BytesToString(recorder.Body.Bytes()),
			recorder.Code)
		/*		if enablePrometheus {
				defaultMetricsRecord(ctxi, r.RequestURI, r.Method, recorder.Code)
			}*/
	}
}
