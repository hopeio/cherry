/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

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
