package plugin

import (
	"github.com/hopeio/cherry/protobuf/utils/enum"
	protogeni "github.com/hopeio/cherry/utils/encoding/protobuf/protogen"
	"google.golang.org/protobuf/compiler/protogen"
)

func TurnOffExtGenAll(f *protogen.File) bool {
	return protogeni.GetOption[bool](f.Desc, enum.E_EnumNoExtgenAll, false)
}

func TurnOffExtGen(e *protogen.Enum) bool {
	return protogeni.GetOption[bool](e.Desc, enum.E_EnumNoExtgen, false)
}

func GetEnumValueCN(ev *protogen.EnumValue) string {
	return protogeni.GetOption[string](ev.Desc, enum.E_EnumvalueCn, "")
}

func GetEnumType(e *protogen.Enum) string {

	return protogeni.GetOption[string](e.Desc, enum.E_EnumCustomtype, "int32")
}

func TurnOffEnumValueMap(e *protogen.Enum) bool {
	return protogeni.GetOption[bool](e.Desc, enum.E_EnumNoGenvaluemap, false)
}

func EnabledEnumNumOrder(e *protogen.Enum) bool {
	return protogeni.GetOption[bool](e.Desc, enum.E_EnumNumorder, false)
}

func EnabledEnumJsonMarshal(f *protogen.File, e *protogen.Enum) bool {
	if protogeni.GetOption[bool](e.Desc, enum.E_EnumJsonmarshal, true) {
		return true
	}
	return protogeni.GetOption[bool](f.Desc, enum.E_EnumJsonmarshalAll, true)
}

func EnabledEnumErrorCode(e *protogen.Enum) bool {
	return protogeni.GetOption[bool](e.Desc, enum.E_EnumErrorcode, false)
}

func EnabledEnumGqlGen(f *protogen.File, e *protogen.Enum) bool {
	if protogeni.GetOption[bool](e.Desc, enum.E_EnumGqlgen, true) {
		return true
	}

	return protogeni.GetOption[bool](f.Desc, enum.E_EnumGqlgenAll, true)
}

func EnabledGoEnumPrefix(f *protogen.File, e *protogen.Enum) bool {
	if protogeni.GetOption[bool](e.Desc, enum.E_EnumNoPrefix, true) {
		return true
	}

	return protogeni.GetOption[bool](f.Desc, enum.E_EnumNoPrefixAll, false)
}

func EnabledEnumStringer(e *protogen.Enum) bool {
	return protogeni.GetOption[bool](e.Desc, enum.E_EnumStringer, true)
}
