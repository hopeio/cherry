# cherry

一个开箱即用，高度集成的微服务组件库,可以快速开发集grpc,http,graphql的云原生微服务

cherry服务器，各种服务接口的保留，集成支持，一个服务暴露grpc,http,graphql接口
- 集成opentelemetry实现调用链路跟踪记录，配合context及utils/log 实现完整的请求链路日志记录
- 集成prometheus及pprof实现性能监控及性能问题排查
- 支持框架自生成的由gin提供支持的grpc转http，也支持原生的grpc-gateway
  ![server](_assets/server.webp)
- 
## quick start
`go get github.com/hopeio/cherry`
### install tools
- `install protoc`[https://github.com/protocolbuffers/protobuf/releases](https://github.com/protocolbuffers/protobuf/releases)
- `go run $(go list -m -f {{.Dir}}  github.com/hopeio/protobuf)/tools/install_tools.go`
### generate protobuf
`protogen go -e -w -v -p _example/proto -g _example/protobuf`
 -e(enum扩展) -w(gin gateway) -q(graphql) -v(生成校验代码) -p proto目录 -g 输出pb.go目录
#### use docker(可选的)
`docker run --rm -v $project:/work jybl/protogen protogen go -e -w -p $proto_path -g $proto_output_path`
### run
`go run _example/user/main.go -c _example/user/config.toml`



```go
package main

import (
	"github.com/hopeio/utils/net/http/gin/handler"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/cherry"
	"github.com/hopeio/initialize"
	"user/protobuf/user"
	uconf "user/confdao"
	udao "user/dao"
	userservice "user/service"
	"github.com/hopeio/utils/log"
	
	"google.golang.org/grpc"
)

func main() {
	//配置初始化应该在第一位
	defer initialize.Start(uconf.Conf, udao.Dao)()
	
  config := uconf.Conf.Server.Origin()
  config.GrpcOptions = []grpc.ServerOption{
    grpc.StatsHandler(otelgrpc.NewServerHandler()),
  }
  cherry.Start(&cherry.Server{
        Config: config,
		GrpcHandler: func(gs *grpc.Server) {
			user.RegisterUserServiceServer(gs, userservice.GetUserService())
		},
		GinHandler: func(app *gin.Engine) {
			_ = user.RegisterUserServiceHandlerServer(app, userservice.GetUserService())
			app.Static("/static", "F:/upload")
		},
        /*	GraphqlHandler: model.NewExecutableSchema(model.Config{
                Resolvers: &model.GQLServer{
                UserService:  service.GetUserService(),
                OauthService: service.GetOauthService(),
            }}),*/
	})
}

```



## TODO
- unit test
- english document
- License


