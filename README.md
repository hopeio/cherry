# cherry

一个开箱即用，高度集成的微服务组件库,可以快速开发集grpc,http,graphql的云原生微服务
## quick start
`go get github.com/hopeio/cherry`
#### install tools
- `install protoc`[https://github.com/protocolbuffers/protobuf/releases](https://github.com/protocolbuffers/protobuf/releases)
- `go run $(go list -m -f {{.Dir}}  github.com/hopeio/protobuf)/tools/install_tools.go`
#### generate protobuf
`protogen go -e -w -v -p $proto_path -g $proto_output_path`
 -e(enum扩展) -w(gin gateway) -q(graphql) -v(生成校验代码) -p proto目录 -g 输出pb.go目录
#### use docker(可选的)
`docker run --rm -v $project:/work jybl/protogen protogen go -e -w -p $proto_path -g $proto_output_path`

##
### initialize
基于反射自动注入的配置及dao注入初始化，并暴露一个全局变量，记录模块信息
![initialize](_readme/assets/initialize.webp)

#### 一个应用的启动，应该如此简单
##### config（配置）
支持nacos,local file,http请求作为配置中心,可扩展支持etcd,apollo,viper(获取配置代理，底层是其他配置中心)，支持viper支持的所有格式("json", "toml", "yaml", "yml", "properties", "props", "prop", "hcl", "tfvars", "dotenv", "env", "ini")的配置文件，
支持dev，test，prod环境本，启动命令区分

##### 配置模板
```toml
# dev | test | stage | prod |...
Env = "dev" # 将会选择与Env名字相同的环境配置
[dev]
ConfigTemplateDir = "." # 模板目录,将会生成配置模板
```
仅需以上最小配置,点击启动,即可生成配置模板
如果还是麻烦,试试直接用 `--format ${配置格式} -e ${环境} -p ${模板路径}`  启动吧
#### 启动配置
仅需配置配置中心,后续配置均从配置中心拉取及自动更新
```toml
Module = "hoper"
# dev | test | stage | prod |...
Env = "dev" # 将会选择与Env名字相同的环境配置

[dev] 
debug = true
ConfigTemplateDir = "." # 将会生成配置模板
# 上方是一个个初始配置,如果不知道如何进行接下来的配置,可以先启动生成配置模板
[dev.ConfigCenter]
Type = "local"
Watch  = true
NoInject = ["Apollo","Etcd", "Es"]

[dev.ConfigCenter.local]
Debug = true
ConfigPath = "local.toml"
ReloadType = "fsnotify"

[dev.ConfigCenter.http]
Interval = 10000000000
Url = "http://localhost:6666/local.toml"

[dev.ConfigCenter.nacos]
DataId = "pro"
Group = "DEFAULT_GROUP"

[[dev.ConfigCenter.nacos.ServerConfigs]]
Scheme = "http"
IpAddr = "nacos"
Port = 9000
GrpcPort = 10000

[dev.ConfigCenter.nacos.ClientConfig]
NamespaceId = "xxx"
username = "nacos"
password = "nacos"
LogLevel = "debug"

```
```go
import(
  "github.com/hopeio/initialize/conf_dao/server"
)
type config struct {
	//自定义的配置
	Customize serverConfig
	Server    server.ServerConfig
}
type serverConfig struct{
    TokenMaxAge time.Duration
}

var Conf = &config{}
// 注入配置前初始化
func (c *config) InitBeforeInject() {
    c.Customize.TokenMaxAge = time.Second * 60 * 60 * 24
}
// 注入配置后初始化
func (c *config) InitAfterInject() {
	c.Customize.TokenMaxAge = time.Second * 60 * 60 * 24 * c.Customize.TokenMaxAge
}

func main() {
    //配置初始化应该在第一位
    defer initialize.Start(Conf, nil)()
}
```
如果还有Dao要初始化
```go
import(
    "github.com/hopeio/initialize/conf_dao/gormdb/postgres"
    initredis "github.com/hopeio/initialize/conf_dao/redis"
)
// dao dao.
type dao struct {
	// GORMDB 数据库连接
	GORMDB   *postgres.DB
	StdDB    *sql.DB
	// RedisPool Redis连接池
	Redis *initredis.Client
}
// 注入配置前初始化
func (c *dao) InitBeforeInject() {
}
// 注入配置后初始化
func (c *dao) InitAfterInjectConfig() {
}
// 注入dao后初始化
func (d *dao) InitAfterInject() {
	db := d.GORMDB
	db.Callback().Create().Remove("gorm:save_before_associations")
	db.Callback().Create().Remove("gorm:save_after_associations")
	db.Callback().Update().Remove("gorm:save_before_associations")
	db.Callback().Update().Remove("gorm:save_after_associations")

	d.StdDB, _ = db.DB()
}
func main() {
    defer initialize.Start(Conf, dao)()
}
```
原生集成了redis,gormdb(mysql,postgressql,sqlite),kafka,pebbledb,apollo,badgerdb,etcd,elasticsearch,nsq,ristretto,viper等，并且非常简单的支持自定义扩展,不局限于Dao对象，任何对象都支持根据配置自动注入生成
## context
一个轻量却强大的上下文管理器,一个请求会生成一个context，贯穿整个请求，context记录原始请求上下文，请求时间，客户端信息，权限校验信息，及负责判断是否内部调用，
及附带唯一traceId的日志记录器
其中权限校验采用jwt，具体的校验模型采用接口，可供使用方自定义
支持http及fasthttp,并支持自定义的请求类型
![context](_readme/assets/context.webp)

### server
cherry服务器，各种服务接口的保留，集成支持，一个服务暴露grpc,http,graphql接口
- 集成opentelemetry实现调用链路跟踪记录，配合context及utils/log 实现完整的请求链路日志记录
- 集成prometheus及pprof实现性能监控及性能问题排查
- 支持框架自生成的由gin提供支持的grpc转http，也支持原生的grpc-gateway
![server](_readme/assets/server.webp)

```go
package main

import (
	"github.com/hopeio/utils/net/http/gin/handler"
	"net/http"
	"time"

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
  server.Start(&server.Server{
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


### utils

各种工具库

以下是一些可以单独成库的工具
#### scheduler
##### engine
并发控制，一个任务调度框架，可以控制goroutine数量,任务失败重试，任务衍生子任务执行，任务检测，任务统计
##### crawler
爬虫框架，基于scheduler/engine

## contribute
### build docker image
```base
`$protobuf_dir/tools/docker/docker_build_local.sh $GOPATH $PROTOC $Image`
```
### upgrade go
`docker build -t jybl/protogen -f $protobuf_dir/tools/docker/Dockerfile_upgrade .`
## TODO
- unit test
- english document
- License


