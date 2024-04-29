// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0-devel
// 	protoc        v5.26.1
// source: cherry/protobuf/request/param.proto

package request

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Page struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PageNo   uint32 `protobuf:"varint,1,opt,name=pageNo,proto3" json:"pageNo,omitempty"`
	PageSize uint32 `protobuf:"varint,2,opt,name=pageSize,proto3" json:"pageSize,omitempty"`
}

func (x *Page) Reset() {
	*x = Page{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cherry_protobuf_request_param_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Page) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Page) ProtoMessage() {}

func (x *Page) ProtoReflect() protoreflect.Message {
	mi := &file_cherry_protobuf_request_param_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Page.ProtoReflect.Descriptor instead.
func (*Page) Descriptor() ([]byte, []int) {
	return file_cherry_protobuf_request_param_proto_rawDescGZIP(), []int{0}
}

func (x *Page) GetPageNo() uint32 {
	if x != nil {
		return x.PageNo
	}
	return 0
}

func (x *Page) GetPageSize() uint32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

type Id struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *Id) Reset() {
	*x = Id{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cherry_protobuf_request_param_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Id) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Id) ProtoMessage() {}

func (x *Id) ProtoReflect() protoreflect.Message {
	mi := &file_cherry_protobuf_request_param_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Id.ProtoReflect.Descriptor instead.
func (*Id) Descriptor() ([]byte, []int) {
	return file_cherry_protobuf_request_param_proto_rawDescGZIP(), []int{1}
}

func (x *Id) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type IdStr struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IdStr string `protobuf:"bytes,1,opt,name=idStr,proto3" json:"idStr,omitempty"`
}

func (x *IdStr) Reset() {
	*x = IdStr{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cherry_protobuf_request_param_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IdStr) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IdStr) ProtoMessage() {}

func (x *IdStr) ProtoReflect() protoreflect.Message {
	mi := &file_cherry_protobuf_request_param_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IdStr.ProtoReflect.Descriptor instead.
func (*IdStr) Descriptor() ([]byte, []int) {
	return file_cherry_protobuf_request_param_proto_rawDescGZIP(), []int{2}
}

func (x *IdStr) GetIdStr() string {
	if x != nil {
		return x.IdStr
	}
	return ""
}

type Cursor struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cursor uint64 `protobuf:"varint,1,opt,name=cursor,proto3" json:"cursor,omitempty"`
}

func (x *Cursor) Reset() {
	*x = Cursor{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cherry_protobuf_request_param_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Cursor) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cursor) ProtoMessage() {}

func (x *Cursor) ProtoReflect() protoreflect.Message {
	mi := &file_cherry_protobuf_request_param_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cursor.ProtoReflect.Descriptor instead.
func (*Cursor) Descriptor() ([]byte, []int) {
	return file_cherry_protobuf_request_param_proto_rawDescGZIP(), []int{3}
}

func (x *Cursor) GetCursor() uint64 {
	if x != nil {
		return x.Cursor
	}
	return 0
}

type CursorStr struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cursor string `protobuf:"bytes,1,opt,name=cursor,proto3" json:"cursor,omitempty"`
}

func (x *CursorStr) Reset() {
	*x = CursorStr{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cherry_protobuf_request_param_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CursorStr) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CursorStr) ProtoMessage() {}

func (x *CursorStr) ProtoReflect() protoreflect.Message {
	mi := &file_cherry_protobuf_request_param_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CursorStr.ProtoReflect.Descriptor instead.
func (*CursorStr) Descriptor() ([]byte, []int) {
	return file_cherry_protobuf_request_param_proto_rawDescGZIP(), []int{4}
}

func (x *CursorStr) GetCursor() string {
	if x != nil {
		return x.Cursor
	}
	return ""
}

var File_cherry_protobuf_request_param_proto protoreflect.FileDescriptor

var file_cherry_protobuf_request_param_proto_rawDesc = []byte{
	0x0a, 0x23, 0x63, 0x68, 0x65, 0x72, 0x72, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x3a,
	0x0a, 0x04, 0x50, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x61, 0x67, 0x65, 0x4e, 0x6f,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x70, 0x61, 0x67, 0x65, 0x4e, 0x6f, 0x12, 0x1a,
	0x0a, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x22, 0x14, 0x0a, 0x02, 0x49, 0x64,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64,
	0x22, 0x1d, 0x0a, 0x05, 0x49, 0x64, 0x53, 0x74, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x64, 0x53,
	0x74, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x69, 0x64, 0x53, 0x74, 0x72, 0x22,
	0x20, 0x0a, 0x06, 0x43, 0x75, 0x72, 0x73, 0x6f, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x75, 0x72,
	0x73, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x63, 0x75, 0x72, 0x73, 0x6f,
	0x72, 0x22, 0x23, 0x0a, 0x09, 0x43, 0x75, 0x72, 0x73, 0x6f, 0x72, 0x53, 0x74, 0x72, 0x12, 0x16,
	0x0a, 0x06, 0x63, 0x75, 0x72, 0x73, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x63, 0x75, 0x72, 0x73, 0x6f, 0x72, 0x42, 0x4e, 0x0a, 0x21, 0x78, 0x79, 0x7a, 0x2e, 0x68, 0x6f,
	0x70, 0x65, 0x72, 0x2e, 0x63, 0x68, 0x65, 0x72, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5a, 0x29, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x6f, 0x70, 0x65, 0x69, 0x6f, 0x2f, 0x63,
	0x68, 0x65, 0x72, 0x72, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cherry_protobuf_request_param_proto_rawDescOnce sync.Once
	file_cherry_protobuf_request_param_proto_rawDescData = file_cherry_protobuf_request_param_proto_rawDesc
)

func file_cherry_protobuf_request_param_proto_rawDescGZIP() []byte {
	file_cherry_protobuf_request_param_proto_rawDescOnce.Do(func() {
		file_cherry_protobuf_request_param_proto_rawDescData = protoimpl.X.CompressGZIP(file_cherry_protobuf_request_param_proto_rawDescData)
	})
	return file_cherry_protobuf_request_param_proto_rawDescData
}

var file_cherry_protobuf_request_param_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_cherry_protobuf_request_param_proto_goTypes = []interface{}{
	(*Page)(nil),      // 0: request.Page
	(*Id)(nil),        // 1: request.Id
	(*IdStr)(nil),     // 2: request.IdStr
	(*Cursor)(nil),    // 3: request.Cursor
	(*CursorStr)(nil), // 4: request.CursorStr
}
var file_cherry_protobuf_request_param_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_cherry_protobuf_request_param_proto_init() }
func file_cherry_protobuf_request_param_proto_init() {
	if File_cherry_protobuf_request_param_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cherry_protobuf_request_param_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Page); i {
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
		file_cherry_protobuf_request_param_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Id); i {
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
		file_cherry_protobuf_request_param_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IdStr); i {
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
		file_cherry_protobuf_request_param_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Cursor); i {
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
		file_cherry_protobuf_request_param_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CursorStr); i {
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
			RawDescriptor: file_cherry_protobuf_request_param_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_cherry_protobuf_request_param_proto_goTypes,
		DependencyIndexes: file_cherry_protobuf_request_param_proto_depIdxs,
		MessageInfos:      file_cherry_protobuf_request_param_proto_msgTypes,
	}.Build()
	File_cherry_protobuf_request_param_proto = out.File
	file_cherry_protobuf_request_param_proto_rawDesc = nil
	file_cherry_protobuf_request_param_proto_goTypes = nil
	file_cherry_protobuf_request_param_proto_depIdxs = nil
}
