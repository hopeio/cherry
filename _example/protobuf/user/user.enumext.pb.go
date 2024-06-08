package user

import (
	errors "errors"
	errorcode "github.com/hopeio/cherry/protobuf/errorcode"
	log "github.com/hopeio/cherry/utils/log"
	strings "github.com/hopeio/cherry/utils/strings"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	strconv "strconv"
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

func (x Gender) MarshalJSON() ([]byte, error) {
	return strings.QuoteToBytes(x.String()), nil
}

func (x *Gender) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '"' {
		value, ok := Gender_value[string(data[1:len(data)-1])]
		if ok {
			*x = Gender(value)
			return nil
		}
	} else {
		value, err := strconv.ParseInt(string(data), 10, 32)
		if err == nil {
			_, ok := Gender_name[int32(value)]
			if ok {
				*x = Gender(value)
				return nil
			}
		}
	}
	return errors.New("invalid enum value: Gender")
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

func (x Role) MarshalJSON() ([]byte, error) {
	return strings.QuoteToBytes(x.String()), nil
}

func (x *Role) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '"' {
		value, ok := Role_value[string(data[1:len(data)-1])]
		if ok {
			*x = Role(value)
			return nil
		}
	} else {
		value, err := strconv.ParseInt(string(data), 10, 32)
		if err == nil {
			_, ok := Role_name[int32(value)]
			if ok {
				*x = Role(value)
				return nil
			}
		}
	}
	return errors.New("invalid enum value: Role")
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

func (x UserStatus) MarshalJSON() ([]byte, error) {
	return strings.QuoteToBytes(x.String()), nil
}

func (x *UserStatus) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '"' {
		value, ok := UserStatus_value[string(data[1:len(data)-1])]
		if ok {
			*x = UserStatus(value)
			return nil
		}
	} else {
		value, err := strconv.ParseInt(string(data), 10, 32)
		if err == nil {
			_, ok := UserStatus_name[int32(value)]
			if ok {
				*x = UserStatus(value)
				return nil
			}
		}
	}
	return errors.New("invalid enum value: UserStatus")
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

func (x UserErr) MarshalJSON() ([]byte, error) {
	return strings.QuoteToBytes(x.String()), nil
}

func (x *UserErr) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '"' {
		value, ok := UserErr_value[string(data[1:len(data)-1])]
		if ok {
			*x = UserErr(value)
			return nil
		}
	} else {
		value, err := strconv.ParseInt(string(data), 10, 32)
		if err == nil {
			_, ok := UserErr_name[int32(value)]
			if ok {
				*x = UserErr(value)
				return nil
			}
		}
	}
	return errors.New("invalid enum value: UserErr")
}

func (x UserErr) Error() string {
	return x.String()
}

func (x UserErr) ErrRep() *errorcode.ErrRep {
	return &errorcode.ErrRep{Code: errorcode.ErrCode(x), Message: x.String()}
}

func (x UserErr) Message(msg string) error {
	return &errorcode.ErrRep{Code: errorcode.ErrCode(x), Message: msg}
}

func (x UserErr) ErrorLog(err error) error {
	log.Error(err)
	return &errorcode.ErrRep{Code: errorcode.ErrCode(x), Message: x.String()}
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
