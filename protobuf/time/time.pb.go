// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v5.26.1
// source: cherry/protobuf/time/time.proto

package time

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

// js不选择支持纳秒级时间戳,都是浮点数,最大53位
type NanoTime struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Nanos int64 `protobuf:"varint,1,opt,name=nanos,proto3" json:"nanos,omitempty"`
}

func (x *NanoTime) Reset() {
	*x = NanoTime{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cherry_protobuf_time_time_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NanoTime) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NanoTime) ProtoMessage() {}

func (x *NanoTime) ProtoReflect() protoreflect.Message {
	mi := &file_cherry_protobuf_time_time_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NanoTime.ProtoReflect.Descriptor instead.
func (*NanoTime) Descriptor() ([]byte, []int) {
	return file_cherry_protobuf_time_time_proto_rawDescGZIP(), []int{0}
}

func (x *NanoTime) GetNanos() int64 {
	if x != nil {
		return x.Nanos
	}
	return 0
}

type MilliTime struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Millis int64 `protobuf:"varint,1,opt,name=millis,proto3" json:"millis,omitempty"`
}

func (x *MilliTime) Reset() {
	*x = MilliTime{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cherry_protobuf_time_time_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MilliTime) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MilliTime) ProtoMessage() {}

func (x *MilliTime) ProtoReflect() protoreflect.Message {
	mi := &file_cherry_protobuf_time_time_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MilliTime.ProtoReflect.Descriptor instead.
func (*MilliTime) Descriptor() ([]byte, []int) {
	return file_cherry_protobuf_time_time_proto_rawDescGZIP(), []int{1}
}

func (x *MilliTime) GetMillis() int64 {
	if x != nil {
		return x.Millis
	}
	return 0
}

type MacroTime struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Macros int64 `protobuf:"varint,1,opt,name=macros,proto3" json:"macros,omitempty"`
}

func (x *MacroTime) Reset() {
	*x = MacroTime{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cherry_protobuf_time_time_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MacroTime) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MacroTime) ProtoMessage() {}

func (x *MacroTime) ProtoReflect() protoreflect.Message {
	mi := &file_cherry_protobuf_time_time_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MacroTime.ProtoReflect.Descriptor instead.
func (*MacroTime) Descriptor() ([]byte, []int) {
	return file_cherry_protobuf_time_time_proto_rawDescGZIP(), []int{2}
}

func (x *MacroTime) GetMacros() int64 {
	if x != nil {
		return x.Macros
	}
	return 0
}

type SecondTime struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Seconds int64 `protobuf:"varint,1,opt,name=seconds,proto3" json:"seconds,omitempty"`
}

func (x *SecondTime) Reset() {
	*x = SecondTime{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cherry_protobuf_time_time_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SecondTime) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SecondTime) ProtoMessage() {}

func (x *SecondTime) ProtoReflect() protoreflect.Message {
	mi := &file_cherry_protobuf_time_time_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SecondTime.ProtoReflect.Descriptor instead.
func (*SecondTime) Descriptor() ([]byte, []int) {
	return file_cherry_protobuf_time_time_proto_rawDescGZIP(), []int{3}
}

func (x *SecondTime) GetSeconds() int64 {
	if x != nil {
		return x.Seconds
	}
	return 0
}

type Date struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Seconds int64 `protobuf:"varint,1,opt,name=seconds,proto3" json:"seconds,omitempty"`
}

func (x *Date) Reset() {
	*x = Date{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cherry_protobuf_time_time_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Date) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Date) ProtoMessage() {}

func (x *Date) ProtoReflect() protoreflect.Message {
	mi := &file_cherry_protobuf_time_time_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Date.ProtoReflect.Descriptor instead.
func (*Date) Descriptor() ([]byte, []int) {
	return file_cherry_protobuf_time_time_proto_rawDescGZIP(), []int{4}
}

func (x *Date) GetSeconds() int64 {
	if x != nil {
		return x.Seconds
	}
	return 0
}

type Duration struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Duration int64 `protobuf:"varint,1,opt,name=duration,proto3" json:"duration,omitempty"`
}

func (x *Duration) Reset() {
	*x = Duration{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cherry_protobuf_time_time_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Duration) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Duration) ProtoMessage() {}

func (x *Duration) ProtoReflect() protoreflect.Message {
	mi := &file_cherry_protobuf_time_time_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Duration.ProtoReflect.Descriptor instead.
func (*Duration) Descriptor() ([]byte, []int) {
	return file_cherry_protobuf_time_time_proto_rawDescGZIP(), []int{5}
}

func (x *Duration) GetDuration() int64 {
	if x != nil {
		return x.Duration
	}
	return 0
}

type Timestamp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Millis int64 `protobuf:"varint,1,opt,name=millis,proto3" json:"millis,omitempty"`
}

func (x *Timestamp) Reset() {
	*x = Timestamp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cherry_protobuf_time_time_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Timestamp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Timestamp) ProtoMessage() {}

func (x *Timestamp) ProtoReflect() protoreflect.Message {
	mi := &file_cherry_protobuf_time_time_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Timestamp.ProtoReflect.Descriptor instead.
func (*Timestamp) Descriptor() ([]byte, []int) {
	return file_cherry_protobuf_time_time_proto_rawDescGZIP(), []int{6}
}

func (x *Timestamp) GetMillis() int64 {
	if x != nil {
		return x.Millis
	}
	return 0
}

type Time struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Represents seconds of UTC time since Unix epoch
	// 1970-01-01T00:00:00Z. Must be from 0001-01-01T00:00:00Z to
	// 9999-12-31T23:59:59Z inclusive.
	Seconds int64 `protobuf:"varint,1,opt,name=seconds,proto3" json:"seconds,omitempty"`
	// Non-negative fractions of a second at nanosecond resolution. Negative
	// second values with fractions must still have non-negative nanos values
	// that count forward in time. Must be from 0 to 999,999,999
	// inclusive.
	Nanos int32 `protobuf:"varint,2,opt,name=nanos,proto3" json:"nanos,omitempty"`
}

func (x *Time) Reset() {
	*x = Time{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cherry_protobuf_time_time_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Time) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Time) ProtoMessage() {}

func (x *Time) ProtoReflect() protoreflect.Message {
	mi := &file_cherry_protobuf_time_time_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Time.ProtoReflect.Descriptor instead.
func (*Time) Descriptor() ([]byte, []int) {
	return file_cherry_protobuf_time_time_proto_rawDescGZIP(), []int{7}
}

func (x *Time) GetSeconds() int64 {
	if x != nil {
		return x.Seconds
	}
	return 0
}

func (x *Time) GetNanos() int32 {
	if x != nil {
		return x.Nanos
	}
	return 0
}

var File_cherry_protobuf_time_time_proto protoreflect.FileDescriptor

var file_cherry_protobuf_time_time_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x63, 0x68, 0x65, 0x72, 0x72, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x22, 0x20, 0x0a, 0x08, 0x4e, 0x61, 0x6e, 0x6f, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x61, 0x6e, 0x6f, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x05, 0x6e, 0x61, 0x6e, 0x6f, 0x73, 0x22, 0x23, 0x0a, 0x09, 0x4d, 0x69, 0x6c,
	0x6c, 0x69, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x69, 0x6c, 0x6c, 0x69, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x6d, 0x69, 0x6c, 0x6c, 0x69, 0x73, 0x22, 0x23,
	0x0a, 0x09, 0x4d, 0x61, 0x63, 0x72, 0x6f, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6d,
	0x61, 0x63, 0x72, 0x6f, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x6d, 0x61, 0x63,
	0x72, 0x6f, 0x73, 0x22, 0x26, 0x0a, 0x0a, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x54, 0x69, 0x6d,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x07, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x22, 0x20, 0x0a, 0x04, 0x44,
	0x61, 0x74, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x22, 0x26, 0x0a,
	0x08, 0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x64, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x23, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x69, 0x6c, 0x6c, 0x69, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x06, 0x6d, 0x69, 0x6c, 0x6c, 0x69, 0x73, 0x22, 0x36, 0x0a, 0x04, 0x54, 0x69,
	0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x07, 0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x12, 0x14, 0x0a, 0x05,
	0x6e, 0x61, 0x6e, 0x6f, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6e, 0x61, 0x6e,
	0x6f, 0x73, 0x42, 0x2a, 0x50, 0x01, 0x5a, 0x26, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x68, 0x6f, 0x70, 0x65, 0x69, 0x6f, 0x2f, 0x63, 0x68, 0x65, 0x72, 0x72, 0x79,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cherry_protobuf_time_time_proto_rawDescOnce sync.Once
	file_cherry_protobuf_time_time_proto_rawDescData = file_cherry_protobuf_time_time_proto_rawDesc
)

func file_cherry_protobuf_time_time_proto_rawDescGZIP() []byte {
	file_cherry_protobuf_time_time_proto_rawDescOnce.Do(func() {
		file_cherry_protobuf_time_time_proto_rawDescData = protoimpl.X.CompressGZIP(file_cherry_protobuf_time_time_proto_rawDescData)
	})
	return file_cherry_protobuf_time_time_proto_rawDescData
}

var file_cherry_protobuf_time_time_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_cherry_protobuf_time_time_proto_goTypes = []interface{}{
	(*NanoTime)(nil),   // 0: time.NanoTime
	(*MilliTime)(nil),  // 1: time.MilliTime
	(*MacroTime)(nil),  // 2: time.MacroTime
	(*SecondTime)(nil), // 3: time.SecondTime
	(*Date)(nil),       // 4: time.Date
	(*Duration)(nil),   // 5: time.Duration
	(*Timestamp)(nil),  // 6: time.Timestamp
	(*Time)(nil),       // 7: time.Time
}
var file_cherry_protobuf_time_time_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_cherry_protobuf_time_time_proto_init() }
func file_cherry_protobuf_time_time_proto_init() {
	if File_cherry_protobuf_time_time_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cherry_protobuf_time_time_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NanoTime); i {
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
		file_cherry_protobuf_time_time_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MilliTime); i {
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
		file_cherry_protobuf_time_time_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MacroTime); i {
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
		file_cherry_protobuf_time_time_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SecondTime); i {
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
		file_cherry_protobuf_time_time_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Date); i {
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
		file_cherry_protobuf_time_time_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Duration); i {
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
		file_cherry_protobuf_time_time_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Timestamp); i {
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
		file_cherry_protobuf_time_time_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Time); i {
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
			RawDescriptor: file_cherry_protobuf_time_time_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_cherry_protobuf_time_time_proto_goTypes,
		DependencyIndexes: file_cherry_protobuf_time_time_proto_depIdxs,
		MessageInfos:      file_cherry_protobuf_time_time_proto_msgTypes,
	}.Build()
	File_cherry_protobuf_time_time_proto = out.File
	file_cherry_protobuf_time_time_proto_rawDesc = nil
	file_cherry_protobuf_time_time_proto_goTypes = nil
	file_cherry_protobuf_time_time_proto_depIdxs = nil
}
