package user

import (
	errors "errors"
)

func (x *User) Validate() error {
	return nil
}
func (x *SignupReq) Validate() error {
	if !(len(x.Password) > 5) {
		return errors.New(`密码最短6位`)
	}
	return nil
}
