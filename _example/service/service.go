/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package service

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/hopeio/cherry/_example/proto"
	"github.com/hopeio/gox/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type UserService struct {
	user.UnimplementedUserServiceServer
}

func (u *UserService) Signup(ctx context.Context, req *pb.SignupReq) (*wrapperspb.StringValue, error) {

	if req.Mail == "" && req.Phone == "" {
		return nil, errors.InvalidArgument.Msg("请填写邮箱或手机号")
	}

	return &wrapperspb.StringValue{Value: "注册成功"}, nil
}

func (u *UserService) GetUser(ctx context.Context, req *pb.GetUserReq) (*pb.User, error) {
	return &user.User{Id: req.Id, ActivatedAt: timestamppb.Now()}, nil
}
