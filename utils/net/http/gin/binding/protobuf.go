// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
)

type protobufBinding struct{}

func (protobufBinding) Name() string {
	return "protobuf"
}

func (b protobufBinding) Bind(ctx *gin.Context, obj interface{}) error {
	buf, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}
	return DecodeProtobuf(buf, obj)
}

func DecodeProtobuf(body []byte, obj interface{}) error {
	if err := proto.Unmarshal(body, obj.(proto.Message)); err != nil {
		return err
	}
	// Here it's same to return Validate(obj), but util now we can't add
	// `binding:""` to the struct which automatically generate by gen-proto
	return nil
	// return Validate(obj)
}
