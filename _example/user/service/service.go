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
	"github.com/hopeio/gox/errors"
	"github.com/hopeio/protobuf/time/timestamp"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/hopeio/cherry/_example/protobuf/user"
	"github.com/hopeio/gox/context/httpctx"
)

type UserService struct {
	user.UnimplementedUserServiceServer
}

func (u *UserService) Signup(ctx context.Context, req *user.SignupReq) (*wrapperspb.StringValue, error) {
	ctxi, _ := httpctx.FromContext(ctx)
	defer ctxi.StartSpanEnd("")()
	if req.Mail == "" && req.Phone == "" {
		return nil, errors.InvalidArgument.Msg("请填写邮箱或手机号")
	}

	return &wrapperspb.StringValue{Value: "注册成功"}, nil
}

func (u *UserService) GetUser(ctx context.Context, req *user.GetUserReq) (*user.User, error) {
	ctxi, _ := httpctx.FromContext(ctx)
	defer ctxi.StartSpanEnd("")()
	return &user.User{Id: req.Id, CreatedAt: timestamp.Now()}, nil
}
func Test(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)
	ctxi, _ := httpctx.FromContext(ctx.Request.Context())
	defer ctxi.StartSpanEnd("")()
	ctx.JSON(200, user.User{Id: uint64(id), CreatedAt: timestamp.Now()})
}
