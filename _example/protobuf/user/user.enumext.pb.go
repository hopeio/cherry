package user

import (
	errors "errors"
	errcode "github.com/hopeio/protobuf/errcode"
	log "github.com/hopeio/utils/log"
	strings "github.com/hopeio/utils/strings"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
)

func (x Gender) String() string {

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
	w.Write(strings.QuoteToBytes(x.String()))
}

func (x *Gender) UnmarshalGQL(v interface{}) error {
	if i, ok := v.(int32); ok {
		*x = Gender(i)
		return nil
	}
	return errors.New("enum need integer type")
}

func (x Role) String() string {

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
	w.Write(strings.QuoteToBytes(x.String()))
}

func (x *Role) UnmarshalGQL(v interface{}) error {
	if i, ok := v.(int32); ok {
		*x = Role(i)
		return nil
	}
	return errors.New("enum need integer type")
}

func (x UserStatus) String() string {

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
	w.Write(strings.QuoteToBytes(x.String()))
}

func (x *UserStatus) UnmarshalGQL(v interface{}) error {
	if i, ok := v.(int32); ok {
		*x = UserStatus(i)
		return nil
	}
	return errors.New("enum need integer type")
}

func (x UserErr) String() string {

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

func (x UserErr) ErrRep() *errcode.ErrRep {
	return &errcode.ErrRep{Code: errcode.ErrCode(x), Message: x.String()}
}

func (x UserErr) Message(msg string) error {
	return &errcode.ErrRep{Code: errcode.ErrCode(x), Message: msg}
}

func (x UserErr) ErrorLog(err error) error {
	log.Error(err)
	return &errcode.ErrRep{Code: errcode.ErrCode(x), Message: x.String()}
}

func (x UserErr) GrpcStatus() *status.Status {
	return status.New(codes.Code(x), x.String())
}

func (x UserErr) MarshalGQL(w io.Writer) {
	w.Write(strings.QuoteToBytes(x.String()))
}

func (x *UserErr) UnmarshalGQL(v interface{}) error {
	if i, ok := v.(int32); ok {
		*x = UserErr(i)
		return nil
	}
	return errors.New("enum need integer type")
}
