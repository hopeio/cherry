// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"github.com/hopeio/cherry/utils/net/http/binding"
	"net/http"

	"github.com/gin-gonic/gin"
)

const defaultMemory = 32 << 20

type formBinding struct{}
type formPostBinding struct{}
type formMultipartBinding struct{}

func (formBinding) Name() string {
	return "form"
}

func (formBinding) Bind(ctx *gin.Context, obj interface{}) error {
	if err := ctx.Request.ParseMultipartForm(defaultMemory); err != nil {
		if err != http.ErrNotMultipart {
			return err
		}
	}
	args := binding.Args{binding.FormSource(ctx.Request.Form)}
	if err := binding.MapForm(obj, args); err != nil {
		return err
	}
	return Validate(obj)
}

func (formPostBinding) Name() string {
	return "form-urlencoded"
}

func (formPostBinding) Bind(ctx *gin.Context, obj interface{}) error {
	if err := ctx.Request.ParseForm(); err != nil {
		return err
	}

	args := binding.Args{binding.FormSource(ctx.Request.Form)}
	if err := binding.MapForm(obj, args); err != nil {
		return err
	}
	return Validate(obj)
}

func (formMultipartBinding) Name() string {
	return "multipart/form-data"
}

func (formMultipartBinding) Bind(ctx *gin.Context, obj interface{}) error {
	if err := ctx.Request.ParseMultipartForm(defaultMemory); err != nil {
		return err
	}
	if err := binding.MappingByPtr(obj, (*binding.MultipartSource)(ctx.Request), Tag); err != nil {
		return err
	}

	return Validate(obj)
}
