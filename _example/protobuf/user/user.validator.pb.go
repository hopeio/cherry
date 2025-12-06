package user

import (
	errors "errors"
	validator "github.com/hopeio/gox/validator"
)

func (x *User) Validate() error {
	if x.CreatedAt != nil {
		if err := validator.ValidateStruct(x.CreatedAt); err != nil {
			return validator.FieldError("CreatedAt", err)
		}
	}
	if x.ActivatedAt != nil {
		if err := validator.ValidateStruct(x.ActivatedAt); err != nil {
			return validator.FieldError("ActivatedAt", err)
		}
	}
	if x.DeletedAt != nil {
		if err := validator.ValidateStruct(x.DeletedAt); err != nil {
			return validator.FieldError("DeletedAt", err)
		}
	}
	return nil
}
func (x *SignupReq) Validate() error {
	if !(len(x.Password) > 5) {
		return errors.New(`密码最短6位`)
	}
	return nil
}
func (x *GetUserReq) Validate() error {
	return nil
}
