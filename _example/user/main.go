package main

import (
	"github.com/hopeio/cherry"
	"github.com/hopeio/cherry/_example/user/api"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"google.golang.org/grpc"
)

func main() {
	conf := cherry.NewConfig()
	conf.GrpcOptions = []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	}

	cherry.Start(&cherry.Server{
		Config: conf,
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
