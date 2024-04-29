package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/cherry/_example/user/dao"
	gormi "github.com/hopeio/cherry/utils/dao/db/gorm"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"strconv"

	"github.com/hopeio/cherry/_example/protobuf/user"
	"github.com/hopeio/cherry/_example/user/confdao"
	"github.com/hopeio/cherry/context/http_context"
	"github.com/hopeio/cherry/protobuf/errorcode"
)

type UserService struct {
	user.UnimplementedUserServiceServer
}

func (u *UserService) Signup(ctx context.Context, req *user.SignupReq) (*wrapperspb.StringValue, error) {
	ctxi, span := http_context.ContextFromContext(ctx).StartSpan("")
	defer span.End()
	ctx = ctxi.Context()
	if req.Mail == "" && req.Phone == "" {
		return nil, errorcode.DBError.Message("请填写邮箱或手机号")
	}

	formatNow := ctxi.TimeString
	var user = &user.User{
		Name: req.Name,

		Mail:   req.Mail,
		Phone:  req.Phone,
		Gender: req.Gender,

		Role:      user.RoleNormal,
		CreatedAt: formatNow,
		Status:    user.UserStatusInActive,
	}

	db := gormi.NewTraceDB(confdao.Dao.GORMDB.DB, ctx, ctxi.TraceID)
	err := db.Create(&user).Error
	if err != nil {
		return nil, ctxi.ErrorLog(errorcode.DBError.Message("新建出错"), err, "UserService.Creat")
	}
	return &wrapperspb.StringValue{Value: "注册成功"}, nil
}

func Test(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)
	ctxi := http_context.ContextFromContext(ctx.Request.Context())
	t, err := dao.GetDao(ctxi, confdao.Dao.GORMDB.DB).GetJsonArrayT(id)
	if err != nil {
		ctx.Writer.WriteString(err.Error())
	}
	ctx.JSON(200, t)
}
