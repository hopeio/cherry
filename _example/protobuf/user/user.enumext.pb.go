package user

import (
	errors "errors"
	strings "github.com/hopeio/cherry/utils/strings"
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

func (x Gender) MarshalJSON() ([]byte, error) {
	return strings.QuoteToBytes(x.String()), nil
}

func (x *Gender) UnmarshalJSON(data []byte) error {
	value, ok := Gender_value[string(data)]
	if ok {
		*x = Gender(value)
		return nil
	}
	return errors.New("无效的Gender")
}

func (x Gender) MarshalGQL(w io.Writer) {
	w.Write(strings.QuoteToBytes(x.String()))
}

func (x *Gender) UnmarshalGQL(v interface{}) error {
	if i, ok := v.(int32); ok {
		*x = Gender(i)
		return nil
	}
	return errors.New("枚举值需要数字类型")
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
	value, ok := Role_value[string(data)]
	if ok {
		*x = Role(value)
		return nil
	}
	return errors.New("无效的Role")
}

func (x Role) MarshalGQL(w io.Writer) {
	w.Write(strings.QuoteToBytes(x.String()))
}

func (x *Role) UnmarshalGQL(v interface{}) error {
	if i, ok := v.(int32); ok {
		*x = Role(i)
		return nil
	}
	return errors.New("枚举值需要数字类型")
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
	value, ok := UserStatus_value[string(data)]
	if ok {
		*x = UserStatus(value)
		return nil
	}
	return errors.New("无效的UserStatus")
}

func (x UserStatus) MarshalGQL(w io.Writer) {
	w.Write(strings.QuoteToBytes(x.String()))
}

func (x *UserStatus) UnmarshalGQL(v interface{}) error {
	if i, ok := v.(int32); ok {
		*x = UserStatus(i)
		return nil
	}
	return errors.New("枚举值需要数字类型")
}
