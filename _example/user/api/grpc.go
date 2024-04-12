package api

import (
	"github.com/hopeio/cherry/_example/protobuf/user"
	userService "github.com/hopeio/cherry/_example/user/service"
	"google.golang.org/grpc"
)

func GrpcRegister(gs *grpc.Server) {
	user.RegisterUserServiceServer(gs, userService.GetUserService())

}
