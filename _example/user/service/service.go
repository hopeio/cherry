/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/utils/errors/errcode"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"strconv"

	"github.com/hopeio/cherry/_example/protobuf/user"
	"github.com/hopeio/context/httpctx"
)

type UserService struct {
	user.UnimplementedUserServiceServer
}

func (u *UserService) Signup(ctx context.Context, req *user.SignupReq) (*wrapperspb.StringValue, error) {
	ctxi, _ := httpctx.FromContextValue(ctx)
	defer ctxi.StartSpanEnd("")()
	if req.Mail == "" && req.Phone == "" {
		return nil, errcode.InvalidArgument.Msg("请填写邮箱或手机号")
	}

	return &wrapperspb.StringValue{Value: "注册成功"}, nil
}

func (u *UserService) GetUser(ctx context.Context, req *user.GetUserReq) (*user.User, error) {
	ctxi, _ := httpctx.FromContextValue(ctx)
	defer ctxi.StartSpanEnd("")()
	return &user.User{Id: req.Id}, nil
}
func Test(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)
	ctxi, _ := httpctx.FromContextValue(ctx.Request.Context())
	defer ctxi.StartSpanEnd("")()
	ctx.JSON(200, user.User{Id: uint64(id)})
}
