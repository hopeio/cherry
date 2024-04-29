package options

import (
	gopb "github.com/hopeio/cherry/protobuf/utils/patch"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

func FileOptions(o *protogen.File) *gopb.FileOptions {
	return proto.GetExtension(o.Desc.Options(), gopb.E_File).(*gopb.FileOptions)
}

func ValueOptions(v *protogen.EnumValue) *gopb.Options {
	return proto.GetExtension(v.Desc.Options(), gopb.E_Value).(*gopb.Options)
}
