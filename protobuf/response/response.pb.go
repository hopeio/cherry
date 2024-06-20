// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v5.26.1
// source: cherry/protobuf/response/response.proto

package response

import (
	_ "github.com/danielvladco/go-proto-gql/pkg/graphqlpb"
	any1 "github.com/hopeio/cherry/protobuf/any"
	_ "github.com/hopeio/cherry/protobuf/utils/patch"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AnyReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    uint32     `protobuf:"varint,1,opt,name=code,proto3" json:"code"`
	Message string     `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Details *anypb.Any `protobuf:"bytes,3,opt,name=details,proto3" json:"details,omitempty"`
}

func (x *AnyReply) Reset() {
	*x = AnyReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cherry_protobuf_response_response_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnyReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnyReply) ProtoMessage() {}

func (x *AnyReply) ProtoReflect() protoreflect.Message {
	mi := &file_cherry_protobuf_response_response_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnyReply.ProtoReflect.Descriptor instead.
func (*AnyReply) Descriptor() ([]byte, []int) {
	return file_cherry_protobuf_response_response_proto_rawDescGZIP(), []int{0}
}

func (x *AnyReply) GetCode() uint32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *AnyReply) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *AnyReply) GetDetails() *anypb.Any {
	if x != nil {
		return x.Details
	}
	return nil
}

type RawReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    uint32 `protobuf:"varint,1,opt,name=code,proto3" json:"code"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	// 字节数组json
	Details *any1.RawJson `protobuf:"bytes,3,opt,name=details,proto3" json:"details,omitempty"`
}

func (x *RawReply) Reset() {
	*x = RawReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cherry_protobuf_response_response_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RawReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RawReply) ProtoMessage() {}

func (x *RawReply) ProtoReflect() protoreflect.Message {
	mi := &file_cherry_protobuf_response_response_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RawReply.ProtoReflect.Descriptor instead.
func (*RawReply) Descriptor() ([]byte, []int) {
	return file_cherry_protobuf_response_response_proto_rawDescGZIP(), []int{1}
}

func (x *RawReply) GetCode() uint32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *RawReply) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *RawReply) GetDetails() *any1.RawJson {
	if x != nil {
		return x.Details
	}
	return nil
}

// 返回数据为字符串，用于新建修改删除类的成功失败提示
type CommonRep struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    uint32 `protobuf:"varint,1,opt,name=code,proto3" json:"code"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Details string `protobuf:"bytes,3,opt,name=details,proto3" json:"details,omitempty"`
}

func (x *CommonRep) Reset() {
	*x = CommonRep{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cherry_protobuf_response_response_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommonRep) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommonRep) ProtoMessage() {}

func (x *CommonRep) ProtoReflect() protoreflect.Message {
	mi := &file_cherry_protobuf_response_response_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommonRep.ProtoReflect.Descriptor instead.
func (*CommonRep) Descriptor() ([]byte, []int) {
	return file_cherry_protobuf_response_response_proto_rawDescGZIP(), []int{2}
}

func (x *CommonRep) GetCode() uint32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *CommonRep) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *CommonRep) GetDetails() string {
	if x != nil {
		return x.Details
	}
	return ""
}

type TinyRep struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    uint32 `protobuf:"varint,1,opt,name=code,proto3" json:"code"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *TinyRep) Reset() {
	*x = TinyRep{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cherry_protobuf_response_response_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TinyRep) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TinyRep) ProtoMessage() {}

func (x *TinyRep) ProtoReflect() protoreflect.Message {
	mi := &file_cherry_protobuf_response_response_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TinyRep.ProtoReflect.Descriptor instead.
func (*TinyRep) Descriptor() ([]byte, []int) {
	return file_cherry_protobuf_response_response_proto_rawDescGZIP(), []int{3}
}

func (x *TinyRep) GetCode() uint32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *TinyRep) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type HttpResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Header []string `protobuf:"bytes,1,rep,name=header,proto3" json:"header,omitempty"`
	Body   []byte   `protobuf:"bytes,2,opt,name=body,proto3" json:"body,omitempty"`
	Status uint32   `protobuf:"varint,3,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *HttpResponse) Reset() {
	*x = HttpResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cherry_protobuf_response_response_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HttpResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HttpResponse) ProtoMessage() {}

func (x *HttpResponse) ProtoReflect() protoreflect.Message {
	mi := &file_cherry_protobuf_response_response_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HttpResponse.ProtoReflect.Descriptor instead.
func (*HttpResponse) Descriptor() ([]byte, []int) {
	return file_cherry_protobuf_response_response_proto_rawDescGZIP(), []int{4}
}

func (x *HttpResponse) GetHeader() []string {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *HttpResponse) GetBody() []byte {
	if x != nil {
		return x.Body
	}
	return nil
}

func (x *HttpResponse) GetStatus() uint32 {
	if x != nil {
		return x.Status
	}
	return 0
}

var File_cherry_protobuf_response_response_proto protoreflect.FileDescriptor

var file_cherry_protobuf_response_response_proto_rawDesc = []byte{
	0x0a, 0x27, 0x63, 0x68, 0x65, 0x72, 0x72, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2f, 0x72, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1d,
	0x63, 0x68, 0x65, 0x72, 0x72, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x61, 0x6e, 0x79, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x24, 0x63,
	0x68, 0x65, 0x72, 0x72, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x75,
	0x74, 0x69, 0x6c, 0x73, 0x2f, 0x70, 0x61, 0x74, 0x63, 0x68, 0x2f, 0x67, 0x6f, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x23, 0x64, 0x61, 0x6e, 0x69, 0x65, 0x6c, 0x76, 0x6c, 0x61, 0x64, 0x63,
	0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x67, 0x72, 0x61, 0x70, 0x68,
	0x71, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x7c, 0x0a, 0x08, 0x41, 0x6e, 0x79, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x12, 0x26, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0d, 0x42, 0x12, 0xd2, 0xb5, 0x03, 0x0e, 0xa2, 0x01, 0x0b, 0x6a, 0x73, 0x6f, 0x6e, 0x3a,
	0x22, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2e, 0x0a, 0x07, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c,
	0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x07, 0x64,
	0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x22, 0x74, 0x0a, 0x08, 0x52, 0x61, 0x77, 0x52, 0x65, 0x70,
	0x6c, 0x79, 0x12, 0x26, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d,
	0x42, 0x12, 0xd2, 0xb5, 0x03, 0x0e, 0xa2, 0x01, 0x0b, 0x6a, 0x73, 0x6f, 0x6e, 0x3a, 0x22, 0x63,
	0x6f, 0x64, 0x65, 0x22, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x26, 0x0a, 0x07, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x61, 0x6e, 0x79, 0x2e, 0x52, 0x61, 0x77, 0x4a,
	0x73, 0x6f, 0x6e, 0x52, 0x07, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x22, 0x67, 0x0a, 0x09,
	0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x52, 0x65, 0x70, 0x12, 0x26, 0x0a, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x12, 0xd2, 0xb5, 0x03, 0x0e, 0xa2, 0x01, 0x0b,
	0x6a, 0x73, 0x6f, 0x6e, 0x3a, 0x22, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x52, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x64,
	0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x73, 0x22, 0x4b, 0x0a, 0x07, 0x54, 0x69, 0x6e, 0x79, 0x52, 0x65, 0x70,
	0x12, 0x26, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x12,
	0xd2, 0xb5, 0x03, 0x0e, 0xa2, 0x01, 0x0b, 0x6a, 0x73, 0x6f, 0x6e, 0x3a, 0x22, 0x63, 0x6f, 0x64,
	0x65, 0x22, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x22, 0x52, 0x0a, 0x0c, 0x48, 0x74, 0x74, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f,
	0x64, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x42, 0x52, 0x0a, 0x22, 0x78, 0x79, 0x7a, 0x2e, 0x68, 0x6f,
	0x70, 0x65, 0x72, 0x2e, 0x63, 0x68, 0x65, 0x72, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x50, 0x01, 0x5a, 0x2a,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x6f, 0x70, 0x65, 0x69,
	0x6f, 0x2f, 0x63, 0x68, 0x65, 0x72, 0x72, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_cherry_protobuf_response_response_proto_rawDescOnce sync.Once
	file_cherry_protobuf_response_response_proto_rawDescData = file_cherry_protobuf_response_response_proto_rawDesc
)

func file_cherry_protobuf_response_response_proto_rawDescGZIP() []byte {
	file_cherry_protobuf_response_response_proto_rawDescOnce.Do(func() {
		file_cherry_protobuf_response_response_proto_rawDescData = protoimpl.X.CompressGZIP(file_cherry_protobuf_response_response_proto_rawDescData)
	})
	return file_cherry_protobuf_response_response_proto_rawDescData
}

var file_cherry_protobuf_response_response_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_cherry_protobuf_response_response_proto_goTypes = []interface{}{
	(*AnyReply)(nil),     // 0: response.AnyReply
	(*RawReply)(nil),     // 1: response.RawReply
	(*CommonRep)(nil),    // 2: response.CommonRep
	(*TinyRep)(nil),      // 3: response.TinyRep
	(*HttpResponse)(nil), // 4: response.HttpResponse
	(*anypb.Any)(nil),    // 5: google.protobuf.Any
	(*any1.RawJson)(nil), // 6: any.RawJson
}
var file_cherry_protobuf_response_response_proto_depIdxs = []int32{
	5, // 0: response.AnyReply.details:type_name -> google.protobuf.Any
	6, // 1: response.RawReply.details:type_name -> any.RawJson
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_cherry_protobuf_response_response_proto_init() }
func file_cherry_protobuf_response_response_proto_init() {
	if File_cherry_protobuf_response_response_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cherry_protobuf_response_response_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AnyReply); i {
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
		file_cherry_protobuf_response_response_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RawReply); i {
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
		file_cherry_protobuf_response_response_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommonRep); i {
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
		file_cherry_protobuf_response_response_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TinyRep); i {
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
		file_cherry_protobuf_response_response_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HttpResponse); i {
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
			RawDescriptor: file_cherry_protobuf_response_response_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_cherry_protobuf_response_response_proto_goTypes,
		DependencyIndexes: file_cherry_protobuf_response_response_proto_depIdxs,
		MessageInfos:      file_cherry_protobuf_response_response_proto_msgTypes,
	}.Build()
	File_cherry_protobuf_response_response_proto = out.File
	file_cherry_protobuf_response_response_proto_rawDesc = nil
	file_cherry_protobuf_response_response_proto_goTypes = nil
	file_cherry_protobuf_response_response_proto_depIdxs = nil
}
