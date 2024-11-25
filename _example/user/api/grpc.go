/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package api

import (
	"github.com/hopeio/cherry/_example/protobuf/user"
	userService "github.com/hopeio/cherry/_example/user/service"
	"google.golang.org/grpc"
)

func GrpcRegister(gs *grpc.Server) {
	user.RegisterUserServiceServer(gs, userService.GetUserService())
}
