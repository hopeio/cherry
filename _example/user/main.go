package main

import (
	"github.com/hopeio/cherry"
	"github.com/hopeio/cherry/_example/user/api"
	"github.com/hopeio/cherry/_example/user/confdao"
	"github.com/hopeio/initialize"
	"github.com/hopeio/initialize/conf_center/nacos"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"google.golang.org/grpc"
)

func main() {
	defer initialize.Start(confdao.Conf, confdao.Dao, nacos.ConfigCenter)()

	config := confdao.Conf.Server.Origin()
	config.GrpcOptions = []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	}

	cherry.Start(&cherry.Server{
		Config: config,
		//为了可以自定义中间件
		GrpcHandler: api.GrpcRegister,
		GinHandler:  api.GinRegister,
		/*		GraphqlHandler: model.NewExecutableSchema(model.Config{
				Resolvers: &model.GQLServer{
					UserService:  service.GetUserService(),
					OauthService: service.GetOauthService(),
				}}),*/
	})
}
