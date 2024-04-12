// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"github.com/gin-gonic/gin"
	"github.com/hopeio/cherry/utils/net/http/request/binding"
)

type queryBinding struct{}

func (queryBinding) Name() string {
	return "query"
}

func (queryBinding) Bind(ctx *gin.Context, obj interface{}) error {
	values := ctx.Request.URL.Query()
	args := binding.Args{binding.FormSource(ctx.Request.Form), binding.FormSource(values), uriSource(ctx.Params)}
	if err := binding.MapForm(obj, args); err != nil {
		return err
	}
	return Validate(obj)
}
