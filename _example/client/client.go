package main

import (
	"context"

	"github.com/hopeio/cherry/_example/protobuf/user"
	"github.com/hopeio/gox/log"
	grpcx "github.com/hopeio/gox/net/http/grpc"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpcx.NewClient("localhost:8080", grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := user.NewUserServiceClient(conn)
	log.Info(client.GetUser(context.Background(), &user.GetUserReq{Id: 1}))
}
