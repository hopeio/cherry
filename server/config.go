package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/cherry/utils/configor"
	"github.com/hopeio/cherry/utils/crypto/tls"
	"github.com/hopeio/cherry/utils/log"
	gini "github.com/hopeio/cherry/utils/net/http/gin"
	"github.com/hopeio/cherry/utils/net/http/grpc/web"
	"github.com/hopeio/cherry/utils/verification/validator"
	"github.com/quic-go/quic-go/http3"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

type Http struct {
	http.Server
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
}
type Http3 struct {
	http3.Server
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
}

type HttpOption struct {
	ExcludeLogPrefixes []string
	IncludeLogPrefixes []string
	StaticFs           []StaticFsConfig `json:"static_fs"`
	Middlewares        []http.HandlerFunc
}

type Config struct {
	ServerName  string
	Http        Http
	Http2       http2.Server
	Http3       *Http3 `json:"http3"`
	StopTimeout time.Duration
	Gin         gini.Config `json:"gin"`
	EnableCors  bool
	Cors        *cors.Options `json:"cors"`
	Middlewares []http.HandlerFunc
	HttpOption  HttpOption
	// Grpc options
	GrpcOptions              []grpc.ServerOption
	EnableGrpcWeb            bool
	GrpcWebOption            []web.Option `json:"grpc_web"`
	Tracing, Metrics, GenDoc bool
	BaseContext              func() context.Context
}

func (c *Config) InitBeforeInject() {
	c.Http.Addr = ":8080"
	c.Http.ReadTimeout = 5 * time.Second
	c.Http.WriteTimeout = 5 * time.Second
	c.StopTimeout = 5 * time.Second
	gin.SetMode(gin.ReleaseMode)
	gin.DisableBindValidation()
	validator.DefaultValidator = nil // 自己做校验
	c.EnableCors = true
	c.Tracing = true
	c.Metrics = true
	c.GenDoc = true
}

func (c *Config) InitAfterInject() {

	if c.Http.Addr == "" {
		c.Http.Addr = ":8080"
	}
	if c.Http3 != nil && c.Http3.Addr == "" {
		c.Http3.Addr = ":8080"
	}
	configor.DurationNotify(c.Http.ReadTimeout, time.Second)
	configor.DurationNotify(c.Http.WriteTimeout, time.Second)
	if c.StopTimeout == 0 {
		c.StopTimeout = 5 * time.Second
	}
	configor.DurationNotify(c.StopTimeout, time.Second)
	if c.Http.CertFile != "" && c.Http.KeyFile != "" {
		tlsConfig, err := tls.Certificate(c.Http.CertFile, c.Http.KeyFile)
		if err != nil {
			log.Fatal(err)
		}
		c.Http.TLSConfig = tlsConfig
	}
	if c.Http3 != nil && c.Http3.CertFile != "" && c.Http3.KeyFile != "" {
		tlsConfig, err := tls.Certificate(c.Http3.CertFile, c.Http3.KeyFile)
		if err != nil {
			log.Fatal(err)
		}
		c.Http3.TLSConfig = tlsConfig
	}

}

func defaultServerConfig() *Config {
	c := &Config{}
	c.InitBeforeInject()
	c.InitAfterInject()
	return c
}

type StaticFsConfig struct {
	Prefix string
	Root   string
}
