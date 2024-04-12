// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0-devel
// 	protoc        v3.20.1
// source: cherry/protobuf/utils/patch/go.proto

package patch

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Options represent Go-specific options for Protobuf messages, fields, oneofs, enums, or enum values.
type Options struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name option renames the generated Go identifier and related identifiers.
	// For a message, this renames the generated Go struct and nested messages or enums, if any.
	// For a message field, this renames the generated Go struct field and getter method.
	// For a oneof field, this renames the generated Go struct field, getter method, interface type, and wrapper types.
	// For an enum, this renames the generated Go type.
	// For an enum value, this renames the generated Go const.
	Name *string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// The getter_name option renames the generated getter method (default: Get<Field>)
	// so a custom getter can be implemented in its place.
	GetterName *string `protobuf:"bytes,10,opt,name=getter_name,json=getterName" json:"getter_name,omitempty"` // TODO: implement this
	// The tags option specifies additional struct tags which are appended a generated Go struct field.
	// This option may be specified on a message field or a oneof field.
	// The value should omit the enclosing backticks.
	Tags *string `protobuf:"bytes,20,opt,name=tags" json:"tags,omitempty"`
	// The stringer_name option renames a generated String() method (if any)
	// so a custom String() method can be implemented in its place.
	StringerName *string `protobuf:"bytes,30,opt,name=stringer_name,json=stringerName" json:"stringer_name,omitempty"` // TODO: implement for messages
}

func (x *Options) Reset() {
	*x = Options{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cherry_protobuf_utils_patch_go_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Options) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Options) ProtoMessage() {}

func (x *Options) ProtoReflect() protoreflect.Message {
	mi := &file_cherry_protobuf_utils_patch_go_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Options.ProtoReflect.Descriptor instead.
func (*Options) Descriptor() ([]byte, []int) {
	return file_cherry_protobuf_utils_patch_go_proto_rawDescGZIP(), []int{0}
}

func (x *Options) GetName() string {
	if x != nil && x.Name != nil {
		return *x.Name
	}
	return ""
}

func (x *Options) GetGetterName() string {
	if x != nil && x.GetterName != nil {
		return *x.GetterName
	}
	return ""
}

func (x *Options) GetTags() string {
	if x != nil && x.Tags != nil {
		return *x.Tags
	}
	return ""
}

func (x *Options) GetStringerName() string {
	if x != nil && x.StringerName != nil {
		return *x.StringerName
	}
	return ""
}

type FileOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NonOmitempty *bool `protobuf:"varint,1,opt,name=non_omitempty,json=nonOmitempty" json:"non_omitempty,omitempty"`
	NoEnumPrefix *bool `protobuf:"varint,10,opt,name=no_enum_prefix,json=noEnumPrefix" json:"no_enum_prefix,omitempty"`
}

func (x *FileOptions) Reset() {
	*x = FileOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cherry_protobuf_utils_patch_go_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileOptions) ProtoMessage() {}

func (x *FileOptions) ProtoReflect() protoreflect.Message {
	mi := &file_cherry_protobuf_utils_patch_go_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileOptions.ProtoReflect.Descriptor instead.
func (*FileOptions) Descriptor() ([]byte, []int) {
	return file_cherry_protobuf_utils_patch_go_proto_rawDescGZIP(), []int{1}
}

func (x *FileOptions) GetNonOmitempty() bool {
	if x != nil && x.NonOmitempty != nil {
		return *x.NonOmitempty
	}
	return false
}

func (x *FileOptions) GetNoEnumPrefix() bool {
	if x != nil && x.NoEnumPrefix != nil {
		return *x.NoEnumPrefix
	}
	return false
}

var file_cherry_protobuf_utils_patch_go_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*Options)(nil),
		Field:         7002,
		Name:          "go.message",
		Tag:           "bytes,7002,opt,name=message",
		Filename:      "cherry/protobuf/utils/patch/go.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*Options)(nil),
		Field:         7002,
		Name:          "go.field",
		Tag:           "bytes,7002,opt,name=field",
		Filename:      "cherry/protobuf/utils/patch/go.proto",
	},
	{
		ExtendedType:  (*descriptorpb.OneofOptions)(nil),
		ExtensionType: (*Options)(nil),
		Field:         7001,
		Name:          "go.oneof",
		Tag:           "bytes,7001,opt,name=oneof",
		Filename:      "cherry/protobuf/utils/patch/go.proto",
	},
	{
		ExtendedType:  (*descriptorpb.EnumOptions)(nil),
		ExtensionType: (*Options)(nil),
		Field:         7002,
		Name:          "go.enum",
		Tag:           "bytes,7002,opt,name=enum",
		Filename:      "cherry/protobuf/utils/patch/go.proto",
	},
	{
		ExtendedType:  (*descriptorpb.EnumOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         7003,
		Name:          "go.cn",
		Tag:           "bytes,7003,opt,name=cn",
		Filename:      "cherry/protobuf/utils/patch/go.proto",
	},
	{
		ExtendedType:  (*descriptorpb.EnumValueOptions)(nil),
		ExtensionType: (*Options)(nil),
		Field:         7002,
		Name:          "go.value",
		Tag:           "bytes,7002,opt,name=value",
		Filename:      "cherry/protobuf/utils/patch/go.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FileOptions)(nil),
		ExtensionType: (*FileOptions)(nil),
		Field:         7002,
		Name:          "go.file",
		Tag:           "bytes,7002,opt,name=file",
		Filename:      "cherry/protobuf/utils/patch/go.proto",
	},
}

// Extension fields to descriptorpb.MessageOptions.
var (
	// optional go.Options message = 7002;
	E_Message = &file_cherry_protobuf_utils_patch_go_proto_extTypes[0]
)

// Extension fields to descriptorpb.FieldOptions.
var (
	// optional go.Options field = 7002;
	E_Field = &file_cherry_protobuf_utils_patch_go_proto_extTypes[1]
)

// Extension fields to descriptorpb.OneofOptions.
var (
	// optional go.Options oneof = 7001;
	E_Oneof = &file_cherry_protobuf_utils_patch_go_proto_extTypes[2]
)

// Extension fields to descriptorpb.EnumOptions.
var (
	// optional go.Options enum = 7002;
	E_Enum = &file_cherry_protobuf_utils_patch_go_proto_extTypes[3]
	// optional string cn = 7003;
	E_Cn = &file_cherry_protobuf_utils_patch_go_proto_extTypes[4]
)

// Extension fields to descriptorpb.EnumValueOptions.
var (
	// optional go.Options value = 7002;
	E_Value = &file_cherry_protobuf_utils_patch_go_proto_extTypes[5]
)

// Extension fields to descriptorpb.FileOptions.
var (
	// optional go.FileOptions file = 7002;
	E_File = &file_cherry_protobuf_utils_patch_go_proto_extTypes[6]
)

var File_cherry_protobuf_utils_patch_go_proto protoreflect.FileDescriptor

var file_cherry_protobuf_utils_patch_go_proto_rawDesc = []byte{
	0x0a, 0x24, 0x63, 0x68, 0x65, 0x72, 0x72, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x75, 0x74, 0x69, 0x6c, 0x73, 0x2f, 0x70, 0x61, 0x74, 0x63, 0x68, 0x2f, 0x67, 0x6f,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x67, 0x6f, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x77, 0x0a, 0x07,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x67,
	0x65, 0x74, 0x74, 0x65, 0x72, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x67, 0x65, 0x74, 0x74, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x61, 0x67, 0x73, 0x18, 0x14, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73,
	0x12, 0x23, 0x0a, 0x0d, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x1e, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x65,
	0x72, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x58, 0x0a, 0x0b, 0x46, 0x69, 0x6c, 0x65, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x12, 0x23, 0x0a, 0x0d, 0x6e, 0x6f, 0x6e, 0x5f, 0x6f, 0x6d, 0x69, 0x74,
	0x65, 0x6d, 0x70, 0x74, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x6e, 0x6f, 0x6e,
	0x4f, 0x6d, 0x69, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x24, 0x0a, 0x0e, 0x6e, 0x6f, 0x5f,
	0x65, 0x6e, 0x75, 0x6d, 0x5f, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x0c, 0x6e, 0x6f, 0x45, 0x6e, 0x75, 0x6d, 0x50, 0x72, 0x65, 0x66, 0x69, 0x78, 0x3a,
	0x47, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xda, 0x36, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x67, 0x6f, 0x2e, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52,
	0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x3a, 0x41, 0x0a, 0x05, 0x66, 0x69, 0x65, 0x6c,
	0x64, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0xda, 0x36, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x67, 0x6f, 0x2e, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x52, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x3a, 0x41, 0x0a, 0x05, 0x6f,
	0x6e, 0x65, 0x6f, 0x66, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4f, 0x6e, 0x65, 0x6f, 0x66, 0x4f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x18, 0xd9, 0x36, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x67, 0x6f, 0x2e,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x05, 0x6f, 0x6e, 0x65, 0x6f, 0x66, 0x3a, 0x3e,
	0x0a, 0x04, 0x65, 0x6e, 0x75, 0x6d, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6e, 0x75, 0x6d, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0xda, 0x36, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x67, 0x6f,
	0x2e, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x04, 0x65, 0x6e, 0x75, 0x6d, 0x3a, 0x2d,
	0x0a, 0x02, 0x63, 0x6e, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6e, 0x75, 0x6d, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x18, 0xdb, 0x36, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x63, 0x6e, 0x3a, 0x45, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x21, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6e, 0x75, 0x6d, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xda, 0x36, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0b, 0x2e, 0x67, 0x6f, 0x2e, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x42, 0x0a, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x1c, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46,
	0x69, 0x6c, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xda, 0x36, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0f, 0x2e, 0x67, 0x6f, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x52, 0x04, 0x66, 0x69, 0x6c, 0x65, 0x42, 0x56, 0x0a, 0x25, 0x78, 0x79, 0x7a, 0x2e,
	0x68, 0x6f, 0x70, 0x65, 0x72, 0x2e, 0x63, 0x68, 0x65, 0x72, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x75, 0x74, 0x69, 0x6c, 0x73, 0x2e, 0x70, 0x61, 0x74, 0x63,
	0x68, 0x5a, 0x2d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x6f,
	0x70, 0x65, 0x69, 0x6f, 0x2f, 0x63, 0x68, 0x65, 0x72, 0x72, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x75, 0x74, 0x69, 0x6c, 0x73, 0x2f, 0x70, 0x61, 0x74, 0x63, 0x68,
}

var (
	file_cherry_protobuf_utils_patch_go_proto_rawDescOnce sync.Once
	file_cherry_protobuf_utils_patch_go_proto_rawDescData = file_cherry_protobuf_utils_patch_go_proto_rawDesc
)

func file_cherry_protobuf_utils_patch_go_proto_rawDescGZIP() []byte {
	file_cherry_protobuf_utils_patch_go_proto_rawDescOnce.Do(func() {
		file_cherry_protobuf_utils_patch_go_proto_rawDescData = protoimpl.X.CompressGZIP(file_cherry_protobuf_utils_patch_go_proto_rawDescData)
	})
	return file_cherry_protobuf_utils_patch_go_proto_rawDescData
}

var file_cherry_protobuf_utils_patch_go_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_cherry_protobuf_utils_patch_go_proto_goTypes = []interface{}{
	(*Options)(nil),                       // 0: go.Options
	(*FileOptions)(nil),                   // 1: go.FileOptions
	(*descriptorpb.MessageOptions)(nil),   // 2: google.protobuf.MessageOptions
	(*descriptorpb.FieldOptions)(nil),     // 3: google.protobuf.FieldOptions
	(*descriptorpb.OneofOptions)(nil),     // 4: google.protobuf.OneofOptions
	(*descriptorpb.EnumOptions)(nil),      // 5: google.protobuf.EnumOptions
	(*descriptorpb.EnumValueOptions)(nil), // 6: google.protobuf.EnumValueOptions
	(*descriptorpb.FileOptions)(nil),      // 7: google.protobuf.FileOptions
}
var file_cherry_protobuf_utils_patch_go_proto_depIdxs = []int32{
	2,  // 0: go.message:extendee -> google.protobuf.MessageOptions
	3,  // 1: go.field:extendee -> google.protobuf.FieldOptions
	4,  // 2: go.oneof:extendee -> google.protobuf.OneofOptions
	5,  // 3: go.enum:extendee -> google.protobuf.EnumOptions
	5,  // 4: go.cn:extendee -> google.protobuf.EnumOptions
	6,  // 5: go.value:extendee -> google.protobuf.EnumValueOptions
	7,  // 6: go.file:extendee -> google.protobuf.FileOptions
	0,  // 7: go.message:type_name -> go.Options
	0,  // 8: go.field:type_name -> go.Options
	0,  // 9: go.oneof:type_name -> go.Options
	0,  // 10: go.enum:type_name -> go.Options
	0,  // 11: go.value:type_name -> go.Options
	1,  // 12: go.file:type_name -> go.FileOptions
	13, // [13:13] is the sub-list for method output_type
	13, // [13:13] is the sub-list for method input_type
	7,  // [7:13] is the sub-list for extension type_name
	0,  // [0:7] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_cherry_protobuf_utils_patch_go_proto_init() }
func file_cherry_protobuf_utils_patch_go_proto_init() {
	if File_cherry_protobuf_utils_patch_go_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cherry_protobuf_utils_patch_go_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Options); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_cherry_protobuf_utils_patch_go_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileOptions); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_cherry_protobuf_utils_patch_go_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 7,
			NumServices:   0,
		},
		GoTypes:           file_cherry_protobuf_utils_patch_go_proto_goTypes,
		DependencyIndexes: file_cherry_protobuf_utils_patch_go_proto_depIdxs,
		MessageInfos:      file_cherry_protobuf_utils_patch_go_proto_msgTypes,
		ExtensionInfos:    file_cherry_protobuf_utils_patch_go_proto_extTypes,
	}.Build()
	File_cherry_protobuf_utils_patch_go_proto = out.File
	file_cherry_protobuf_utils_patch_go_proto_rawDesc = nil
	file_cherry_protobuf_utils_patch_go_proto_goTypes = nil
	file_cherry_protobuf_utils_patch_go_proto_depIdxs = nil
}
