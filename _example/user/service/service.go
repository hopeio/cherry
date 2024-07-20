package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"strconv"

	"github.com/hopeio/cherry/_example/protobuf/user"
	"github.com/hopeio/context/httpctx"
	"github.com/hopeio/protobuf/errcode"
)

type UserService struct {
	user.UnimplementedUserServiceServer
}

func (u *UserService) Signup(ctx context.Context, req *user.SignupReq) (*wrapperspb.StringValue, error) {
	ctxi := httpctx.FromContextValue(ctx)
	defer ctxi.StartSpanEnd("")()
	if req.Mail == "" && req.Phone == "" {
		return nil, errcode.DBError.Message("请填写邮箱或手机号")
	}

	return &wrapperspb.StringValue{Value: "注册成功"}, nil
}

func Test(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)
	ctxi := httpctx.FromContextValue(ctx.Request.Context())
	defer ctxi.StartSpanEnd("")()
	ctx.JSON(200, id)
}
