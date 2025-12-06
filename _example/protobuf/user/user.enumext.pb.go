package user

import (
	errors "errors"
	errors1 "github.com/hopeio/gox/errors"
	grpc "github.com/hopeio/gox/net/http/grpc"
	strings "github.com/hopeio/gox/strings"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
)

func (x Gender) Comment() string {
	switch x {
	case GenderPlaceholder:
		return "占位"
	case GenderUnfilled:
		return "未填"
	case GenderMale:
		return "男"
	case GenderFemale:
		return "女"
	}
	return ""
}

func (x Gender) MarshalGQL(w io.Writer) {
	w.Write(strings.SimpleQuoteToBytes(x.String()))
}

func (x *Gender) UnmarshalGQL(v interface{}) error {
	if i, ok := v.(int32); ok {
		*x = Gender(i)
		return nil
	}
	return errors.New("enum need integer type")
}

func (x Role) Comment() string {
	switch x {
	case PlaceholderRole:
		return "占位"
	case RoleNormal:
		return "普通用户"
	case RoleAdmin:
		return "管理员"
	case RoleSuperAdmin:
		return "超级管理员"
	}
	return ""
}

func (x Role) MarshalGQL(w io.Writer) {
	w.Write(strings.SimpleQuoteToBytes(x.String()))
}

func (x *Role) UnmarshalGQL(v interface{}) error {
	if i, ok := v.(int32); ok {
		*x = Role(i)
		return nil
	}
	return errors.New("enum need integer type")
}

func (x UserStatus) Comment() string {
	switch x {
	case UserStatusPlaceholder:
		return "占位"
	case UserStatusInActive:
		return "未激活"
	case UserStatusActivated:
		return "已激活"
	case UserStatusFrozen:
		return "已冻结"
	case UserStatusDeleted:
		return "已注销"
	}
	return ""
}

func (x UserStatus) MarshalGQL(w io.Writer) {
	w.Write(strings.SimpleQuoteToBytes(x.String()))
}

func (x *UserStatus) UnmarshalGQL(v interface{}) error {
	if i, ok := v.(int32); ok {
		*x = UserStatus(i)
		return nil
	}
	return errors.New("enum need integer type")
}

func (x UserErr) Comment() string {
	switch x {
	case UserErrPlaceholder:
		return "占位"
	case UserErrLogin:
		return "用户名或密码错误"
	case UserErrNoActive:
		return "未激活账号"
	case UserErrNoAuthority:
		return "无权限"
	case UserErrLoginTimeout:
		return "登录超时"
	case UserErrInvalidToken:
		return "Token错误"
	case UserErrNoLogin:
		return "未登录"
	}
	return ""
}

func (x UserErr) Error() string {
	return x.String()
}

func (x UserErr) ErrResp() *grpc.ErrResp {
	return &grpc.ErrResp{Code: errors1.ErrCode(x), Msg: x.String()}
}

func (x UserErr) Msg(msg string) *grpc.ErrResp {
	return &grpc.ErrResp{Code: errors1.ErrCode(x), Msg: msg}
}

func (x UserErr) Wrap(err error) *grpc.ErrResp {
	return &grpc.ErrResp{Code: errors1.ErrCode(x), Msg: err.Error()}
}

func (x UserErr) GRPCStatus() *status.Status {
	return status.New(codes.Code(x), x.String())
}

func init() {
	for code, msg := range UserErr_name {
		errors1.Register(errors1.ErrCode(code), msg)
	}
}

func (x UserErr) MarshalGQL(w io.Writer) {
	w.Write(strings.SimpleQuoteToBytes(x.String()))
}

func (x *UserErr) UnmarshalGQL(v interface{}) error {
	if i, ok := v.(int32); ok {
		*x = UserErr(i)
		return nil
	}
	return errors.New("enum need integer type")
}
