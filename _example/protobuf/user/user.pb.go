// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v3.20.1
// source: user/user.proto

package user

import (
	_ "github.com/danielvladco/go-proto-gql/pkg/graphqlpb"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "github.com/hopeio/cherry/protobuf/utils/enum"
	_ "github.com/hopeio/cherry/protobuf/utils/patch"
	_ "github.com/hopeio/cherry/protobuf/utils/validator"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 用户性别
type Gender int32

const (
	GenderPlaceholder Gender = 0
	GenderUnfilled    Gender = 1
	GenderMale        Gender = 2
	GenderFemale      Gender = 3
)

// Enum value maps for Gender.
var (
	Gender_name = map[int32]string{
		0: "GenderPlaceholder",
		1: "GenderUnfilled",
		2: "GenderMale",
		3: "GenderFemale",
	}
	Gender_value = map[string]int32{
		"GenderPlaceholder": 0,
		"GenderUnfilled":    1,
		"GenderMale":        2,
		"GenderFemale":      3,
	}
)

func (x Gender) Enum() *Gender {
	p := new(Gender)
	*p = x
	return p
}

func (x Gender) OrigString() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Gender) Descriptor() protoreflect.EnumDescriptor {
	return file_user_user_proto_enumTypes[0].Descriptor()
}

func (Gender) Type() protoreflect.EnumType {
	return &file_user_user_proto_enumTypes[0]
}

func (x Gender) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Gender.Descriptor instead.
func (Gender) EnumDescriptor() ([]byte, []int) {
	return file_user_user_proto_rawDescGZIP(), []int{0}
}

// 用户角色
type Role int32

const (
	PlaceholderRole Role = 0
	RoleNormal      Role = 1
	RoleAdmin       Role = 2
	RoleSuperAdmin  Role = 3
)

// Enum value maps for Role.
var (
	Role_name = map[int32]string{
		0: "PlaceholderRole",
		1: "RoleNormal",
		2: "RoleAdmin",
		3: "RoleSuperAdmin",
	}
	Role_value = map[string]int32{
		"PlaceholderRole": 0,
		"RoleNormal":      1,
		"RoleAdmin":       2,
		"RoleSuperAdmin":  3,
	}
)

func (x Role) Enum() *Role {
	p := new(Role)
	*p = x
	return p
}

func (x Role) OrigString() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Role) Descriptor() protoreflect.EnumDescriptor {
	return file_user_user_proto_enumTypes[1].Descriptor()
}

func (Role) Type() protoreflect.EnumType {
	return &file_user_user_proto_enumTypes[1]
}

func (x Role) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Role.Descriptor instead.
func (Role) EnumDescriptor() ([]byte, []int) {
	return file_user_user_proto_rawDescGZIP(), []int{1}
}

// 用户角色
type UserStatus int32

const (
	UserStatusPlaceholder UserStatus = 0
	UserStatusInActive    UserStatus = 1
	UserStatusActivated   UserStatus = 2
	UserStatusFrozen      UserStatus = 3
	UserStatusDeleted     UserStatus = 4
)

// Enum value maps for UserStatus.
var (
	UserStatus_name = map[int32]string{
		0: "UserStatusPlaceholder",
		1: "UserStatusInActive",
		2: "UserStatusActivated",
		3: "UserStatusFrozen",
		4: "UserStatusDeleted",
	}
	UserStatus_value = map[string]int32{
		"UserStatusPlaceholder": 0,
		"UserStatusInActive":    1,
		"UserStatusActivated":   2,
		"UserStatusFrozen":      3,
		"UserStatusDeleted":     4,
	}
)

func (x UserStatus) Enum() *UserStatus {
	p := new(UserStatus)
	*p = x
	return p
}

func (x UserStatus) OrigString() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UserStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_user_user_proto_enumTypes[2].Descriptor()
}

func (UserStatus) Type() protoreflect.EnumType {
	return &file_user_user_proto_enumTypes[2]
}

func (x UserStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UserStatus.Descriptor instead.
func (UserStatus) EnumDescriptor() ([]byte, []int) {
	return file_user_user_proto_rawDescGZIP(), []int{2}
}

type UserErr int32

const (
	UserErrPlaceholder  UserErr = 0
	UserErrLogin        UserErr = 1000
	UserErrNoActive     UserErr = 1001
	UserErrNoAuthority  UserErr = 1002
	UserErrLoginTimeout UserErr = 1003
	UserErrInvalidToken UserErr = 1004
	UserErrNoLogin      UserErr = 1005
)

// Enum value maps for UserErr.
var (
	UserErr_name = map[int32]string{
		0:    "UserErrPlaceholder",
		1000: "UserErrLogin",
		1001: "UserErrNoActive",
		1002: "UserErrNoAuthority",
		1003: "UserErrLoginTimeout",
		1004: "UserErrInvalidToken",
		1005: "UserErrNoLogin",
	}
	UserErr_value = map[string]int32{
		"UserErrPlaceholder":  0,
		"UserErrLogin":        1000,
		"UserErrNoActive":     1001,
		"UserErrNoAuthority":  1002,
		"UserErrLoginTimeout": 1003,
		"UserErrInvalidToken": 1004,
		"UserErrNoLogin":      1005,
	}
)

func (x UserErr) Enum() *UserErr {
	p := new(UserErr)
	*p = x
	return p
}

func (x UserErr) OrigString() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (UserErr) Descriptor() protoreflect.EnumDescriptor {
	return file_user_user_proto_enumTypes[3].Descriptor()
}

func (UserErr) Type() protoreflect.EnumType {
	return &file_user_user_proto_enumTypes[3]
}

func (x UserErr) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use UserErr.Descriptor instead.
func (UserErr) EnumDescriptor() ([]byte, []int) {
	return file_user_user_proto_rawDescGZIP(), []int{3}
}

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" gorm:"primaryKey;"`
	Name     string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty" gorm:"size:10;not null" comment:"昵称"`
	Password string `protobuf:"bytes,5,opt,name=password,proto3" json:"-" gorm:"size:32;not null" validate:"gte=8,lte=15" comment:"密码"`
	Mail     string `protobuf:"bytes,6,opt,name=mail,proto3" json:"mail,omitempty" gorm:"size:32" validate:"email" comment:"邮箱"`
	Phone    string `protobuf:"bytes,7,opt,name=phone,proto3" json:"phone,omitempty" gorm:"size:32" validate:"phone" comment:"手机号"`
	// 性别，0未填写，1男，2女
	Gender      Gender     `protobuf:"varint,8,opt,name=gender,proto3,enum=user.Gender" json:"gender,omitempty" gorm:"type:int2;default:0"`
	Role        Role       `protobuf:"varint,24,opt,name=role,proto3,enum=user.Role" json:"role,omitempty" gorm:"type:int2;default:0"`
	Status      UserStatus `protobuf:"varint,28,opt,name=status,proto3,enum=user.UserStatus" json:"status,omitempty" gorm:"type:int2;default:0"`
	CreatedAt   string     `protobuf:"bytes,25,opt,name=createdAt,proto3" json:"createdAt,omitempty" gorm:"type:timestamptz(6);default:now();index"`
	ActivatedAt string     `protobuf:"bytes,3,opt,name=activatedAt,proto3" json:"activatedAt,omitempty" gorm:"<-:false;type:timestamptz(6);index"`
	DeletedAt   string     `protobuf:"bytes,27,opt,name=deletedAt,proto3" json:"deletedAt,omitempty" gorm:"<-:false;type:timestamptz(6);index"` // uint32 isDeleted = 29 [(go.field) = {tags:'gorm:"default:0"'}];
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_user_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_user_user_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *User) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *User) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *User) GetMail() string {
	if x != nil {
		return x.Mail
	}
	return ""
}

func (x *User) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *User) GetGender() Gender {
	if x != nil {
		return x.Gender
	}
	return GenderPlaceholder
}

func (x *User) GetRole() Role {
	if x != nil {
		return x.Role
	}
	return PlaceholderRole
}

func (x *User) GetStatus() UserStatus {
	if x != nil {
		return x.Status
	}
	return UserStatusPlaceholder
}

func (x *User) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

func (x *User) GetActivatedAt() string {
	if x != nil {
		return x.ActivatedAt
	}
	return ""
}

func (x *User) GetDeletedAt() string {
	if x != nil {
		return x.DeletedAt
	}
	return ""
}

type SignupReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 密码
	Password string `protobuf:"bytes,4,opt,name=password,proto3" json:"password,omitempty" validate:"required,gte=6,lte=15" comment:"密码"`
	Name     string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty" validate:"required,gte=3,lte=10" comment:"昵称"`
	Gender   Gender `protobuf:"varint,3,opt,name=gender,proto3,enum=user.Gender" json:"gender,omitempty" validate:"required" comment:"性别"`
	// 邮箱
	Mail string `protobuf:"bytes,6,opt,name=mail,proto3" json:"mail,omitempty" validate:"omitempty,email" comment:"邮箱"`
	// 手机号
	Phone string `protobuf:"bytes,7,opt,name=phone,proto3" json:"phone,omitempty" validate:"omitempty,phone" comment:"手机号"`
	// 验证码
	VCode string `protobuf:"bytes,8,opt,name=vCode,proto3" json:"vCode,omitempty" validate:"required" comment:"验证码"`
}

func (x *SignupReq) Reset() {
	*x = SignupReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignupReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignupReq) ProtoMessage() {}

func (x *SignupReq) ProtoReflect() protoreflect.Message {
	mi := &file_user_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignupReq.ProtoReflect.Descriptor instead.
func (*SignupReq) Descriptor() ([]byte, []int) {
	return file_user_user_proto_rawDescGZIP(), []int{1}
}

func (x *SignupReq) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *SignupReq) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SignupReq) GetGender() Gender {
	if x != nil {
		return x.Gender
	}
	return GenderPlaceholder
}

func (x *SignupReq) GetMail() string {
	if x != nil {
		return x.Mail
	}
	return ""
}

func (x *SignupReq) GetPhone() string {
	if x != nil {
		return x.Phone
	}
	return ""
}

func (x *SignupReq) GetVCode() string {
	if x != nil {
		return x.VCode
	}
	return ""
}

var File_user_user_proto protoreflect.FileDescriptor

var file_user_user_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x04, 0x75, 0x73, 0x65, 0x72, 0x1a, 0x25, 0x63, 0x68, 0x65, 0x72, 0x72, 0x79, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x75, 0x74, 0x69, 0x6c, 0x73, 0x2f, 0x65,
	0x6e, 0x75, 0x6d, 0x2f, 0x65, 0x6e, 0x75, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x24,
	0x63, 0x68, 0x65, 0x72, 0x72, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x75, 0x74, 0x69, 0x6c, 0x73, 0x2f, 0x70, 0x61, 0x74, 0x63, 0x68, 0x2f, 0x67, 0x6f, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e,
	0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2f, 0x63, 0x68, 0x65, 0x72, 0x72, 0x79, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x75, 0x74, 0x69, 0x6c, 0x73, 0x2f, 0x76, 0x61, 0x6c, 0x69,
	0x64, 0x61, 0x74, 0x6f, 0x72, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x6f, 0x72, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x23, 0x64, 0x61, 0x6e, 0x69, 0x65, 0x6c, 0x76, 0x6c, 0x61, 0x64, 0x63,
	0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x67, 0x72, 0x61, 0x70, 0x68,
	0x71, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe1, 0x06, 0x0a, 0x04, 0x55, 0x73, 0x65,
	0x72, 0x12, 0x29, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x42, 0x19, 0xd2,
	0xb5, 0x03, 0x15, 0xa2, 0x01, 0x12, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x70, 0x72, 0x69, 0x6d,
	0x61, 0x72, 0x79, 0x4b, 0x65, 0x79, 0x3b, 0x22, 0x52, 0x02, 0x69, 0x64, 0x12, 0x43, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x2f, 0xd2, 0xb5, 0x03, 0x2b,
	0xa2, 0x01, 0x28, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x73, 0x69, 0x7a, 0x65, 0x3a, 0x31, 0x30,
	0x3b, 0x6e, 0x6f, 0x74, 0x20, 0x6e, 0x75, 0x6c, 0x6c, 0x22, 0x20, 0x63, 0x6f, 0x6d, 0x6d, 0x65,
	0x6e, 0x74, 0x3a, 0x22, 0xe6, 0x98, 0xb5, 0xe7, 0xa7, 0xb0, 0x22, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x6c, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x50, 0xd2, 0xb5, 0x03, 0x4c, 0xa2, 0x01, 0x49, 0x6a, 0x73, 0x6f, 0x6e,
	0x3a, 0x22, 0x2d, 0x22, 0x20, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x73, 0x69, 0x7a, 0x65, 0x3a,
	0x33, 0x32, 0x3b, 0x6e, 0x6f, 0x74, 0x20, 0x6e, 0x75, 0x6c, 0x6c, 0x22, 0x20, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x3a, 0x22, 0x67, 0x74, 0x65, 0x3d, 0x38, 0x2c, 0x6c, 0x74, 0x65,
	0x3d, 0x31, 0x35, 0x22, 0x20, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x3a, 0x22, 0xe5, 0xaf,
	0x86, 0xe7, 0xa0, 0x81, 0x22, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12,
	0x4b, 0x0a, 0x04, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x42, 0x37, 0xd2,
	0xb5, 0x03, 0x33, 0xa2, 0x01, 0x30, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x73, 0x69, 0x7a, 0x65,
	0x3a, 0x33, 0x32, 0x22, 0x20, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x3a, 0x22, 0x65,
	0x6d, 0x61, 0x69, 0x6c, 0x22, 0x20, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x3a, 0x22, 0xe9,
	0x82, 0xae, 0xe7, 0xae, 0xb1, 0x22, 0x52, 0x04, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x50, 0x0a, 0x05,
	0x70, 0x68, 0x6f, 0x6e, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x42, 0x3a, 0xd2, 0xb5, 0x03,
	0x36, 0xa2, 0x01, 0x33, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x73, 0x69, 0x7a, 0x65, 0x3a, 0x33,
	0x32, 0x22, 0x20, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x3a, 0x22, 0x70, 0x68, 0x6f,
	0x6e, 0x65, 0x22, 0x20, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x3a, 0x22, 0xe6, 0x89, 0x8b,
	0xe6, 0x9c, 0xba, 0xe5, 0x8f, 0xb7, 0x22, 0x52, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x12, 0x4e,
	0x0a, 0x06, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c,
	0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x42, 0x28, 0x92, 0x41,
	0x04, 0x9a, 0x02, 0x01, 0x03, 0xd2, 0xb5, 0x03, 0x1d, 0xa2, 0x01, 0x1a, 0x67, 0x6f, 0x72, 0x6d,
	0x3a, 0x22, 0x74, 0x79, 0x70, 0x65, 0x3a, 0x69, 0x6e, 0x74, 0x32, 0x3b, 0x64, 0x65, 0x66, 0x61,
	0x75, 0x6c, 0x74, 0x3a, 0x30, 0x22, 0x52, 0x06, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x41,
	0x0a, 0x04, 0x72, 0x6f, 0x6c, 0x65, 0x18, 0x18, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0a, 0x2e, 0x75,
	0x73, 0x65, 0x72, 0x2e, 0x52, 0x6f, 0x6c, 0x65, 0x42, 0x21, 0xd2, 0xb5, 0x03, 0x1d, 0xa2, 0x01,
	0x1a, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x74, 0x79, 0x70, 0x65, 0x3a, 0x69, 0x6e, 0x74, 0x32,
	0x3b, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x3a, 0x30, 0x22, 0x52, 0x04, 0x72, 0x6f, 0x6c,
	0x65, 0x12, 0x50, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x1c, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x10, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x42, 0x26, 0x92, 0x41, 0x02, 0x40, 0x01, 0xd2, 0xb5, 0x03, 0x1d, 0xa2, 0x01,
	0x1a, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x74, 0x79, 0x70, 0x65, 0x3a, 0x69, 0x6e, 0x74, 0x32,
	0x3b, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x3a, 0x30, 0x22, 0x52, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x53, 0x0a, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x18, 0x19, 0x20, 0x01, 0x28, 0x09, 0x42, 0x35, 0xd2, 0xb5, 0x03, 0x31, 0xa2, 0x01, 0x2e, 0x67,
	0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x74, 0x79, 0x70, 0x65, 0x3a, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x74, 0x7a, 0x28, 0x36, 0x29, 0x3b, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74,
	0x3a, 0x6e, 0x6f, 0x77, 0x28, 0x29, 0x3b, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x22, 0x52, 0x09, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x52, 0x0a, 0x0b, 0x61, 0x63, 0x74, 0x69,
	0x76, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x30, 0xd2,
	0xb5, 0x03, 0x2c, 0xa2, 0x01, 0x29, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x3c, 0x2d, 0x3a, 0x66,
	0x61, 0x6c, 0x73, 0x65, 0x3b, 0x74, 0x79, 0x70, 0x65, 0x3a, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x74, 0x7a, 0x28, 0x36, 0x29, 0x3b, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x22, 0x52,
	0x0b, 0x61, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x4e, 0x0a, 0x09,
	0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x1b, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x30, 0xd2, 0xb5, 0x03, 0x2c, 0xa2, 0x01, 0x29, 0x67, 0x6f, 0x72, 0x6d, 0x3a, 0x22, 0x3c, 0x2d,
	0x3a, 0x66, 0x61, 0x6c, 0x73, 0x65, 0x3b, 0x74, 0x79, 0x70, 0x65, 0x3a, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x74, 0x7a, 0x28, 0x36, 0x29, 0x3b, 0x69, 0x6e, 0x64, 0x65, 0x78,
	0x22, 0x52, 0x09, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x83, 0x04, 0x0a,
	0x09, 0x53, 0x69, 0x67, 0x6e, 0x75, 0x70, 0x52, 0x65, 0x71, 0x12, 0x7a, 0x0a, 0x08, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x42, 0x5e, 0x92, 0x41,
	0x0b, 0x2a, 0x06, 0xe5, 0xaf, 0x86, 0xe7, 0xa0, 0x81, 0x80, 0x01, 0x06, 0xd2, 0xb5, 0x03, 0x34,
	0xa2, 0x01, 0x31, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x3a, 0x22, 0x72, 0x65, 0x71,
	0x75, 0x69, 0x72, 0x65, 0x64, 0x2c, 0x67, 0x74, 0x65, 0x3d, 0x36, 0x2c, 0x6c, 0x74, 0x65, 0x3d,
	0x31, 0x35, 0x22, 0x20, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x3a, 0x22, 0xe5, 0xaf, 0x86,
	0xe7, 0xa0, 0x81, 0x22, 0xe2, 0xdf, 0x1f, 0x14, 0x2a, 0x10, 0xe5, 0xaf, 0x86, 0xe7, 0xa0, 0x81,
	0xe6, 0x9c, 0x80, 0xe7, 0x9f, 0xad, 0x36, 0xe4, 0xbd, 0x8d, 0x70, 0x05, 0x52, 0x08, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x4c, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x38, 0xd2, 0xb5, 0x03, 0x34, 0xa2, 0x01, 0x31, 0x76, 0x61, 0x6c,
	0x69, 0x64, 0x61, 0x74, 0x65, 0x3a, 0x22, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x2c,
	0x67, 0x74, 0x65, 0x3d, 0x33, 0x2c, 0x6c, 0x74, 0x65, 0x3d, 0x31, 0x30, 0x22, 0x20, 0x63, 0x6f,
	0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x3a, 0x22, 0xe6, 0x98, 0xb5, 0xe7, 0xa7, 0xb0, 0x22, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x51, 0x0a, 0x06, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x47, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x42, 0x2b, 0xd2, 0xb5, 0x03, 0x27, 0xa2, 0x01, 0x24, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x3a, 0x22, 0x72, 0x65, 0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x22, 0x20, 0x63,
	0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x3a, 0x22, 0xe6, 0x80, 0xa7, 0xe5, 0x88, 0xab, 0x22, 0x52,
	0x06, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x46, 0x0a, 0x04, 0x6d, 0x61, 0x69, 0x6c, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x42, 0x32, 0xd2, 0xb5, 0x03, 0x2e, 0xa2, 0x01, 0x2b, 0x76, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x3a, 0x22, 0x6f, 0x6d, 0x69, 0x74, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2c, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0x20, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74,
	0x3a, 0x22, 0xe9, 0x82, 0xae, 0xe7, 0xae, 0xb1, 0x22, 0x52, 0x04, 0x6d, 0x61, 0x69, 0x6c, 0x12,
	0x4b, 0x0a, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x42, 0x35,
	0xd2, 0xb5, 0x03, 0x31, 0xa2, 0x01, 0x2e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x3a,
	0x22, 0x6f, 0x6d, 0x69, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2c, 0x70, 0x68, 0x6f, 0x6e, 0x65,
	0x22, 0x20, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x3a, 0x22, 0xe6, 0x89, 0x8b, 0xe6, 0x9c,
	0xba, 0xe5, 0x8f, 0xb7, 0x22, 0x52, 0x05, 0x70, 0x68, 0x6f, 0x6e, 0x65, 0x12, 0x44, 0x0a, 0x05,
	0x76, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x42, 0x2e, 0xd2, 0xb5, 0x03,
	0x2a, 0xa2, 0x01, 0x27, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x3a, 0x22, 0x72, 0x65,
	0x71, 0x75, 0x69, 0x72, 0x65, 0x64, 0x22, 0x20, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x3a,
	0x22, 0xe9, 0xaa, 0x8c, 0xe8, 0xaf, 0x81, 0xe7, 0xa0, 0x81, 0x22, 0x52, 0x05, 0x76, 0x43, 0x6f,
	0x64, 0x65, 0x2a, 0x92, 0x01, 0x0a, 0x06, 0x47, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x21, 0x0a,
	0x11, 0x47, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x68, 0x6f, 0x6c, 0x64,
	0x65, 0x72, 0x10, 0x00, 0x1a, 0x0a, 0x92, 0x9d, 0x20, 0x06, 0xe5, 0x8d, 0xa0, 0xe4, 0xbd, 0x8d,
	0x12, 0x1e, 0x0a, 0x0e, 0x47, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x55, 0x6e, 0x66, 0x69, 0x6c, 0x6c,
	0x65, 0x64, 0x10, 0x01, 0x1a, 0x0a, 0x92, 0x9d, 0x20, 0x06, 0xe6, 0x9c, 0xaa, 0xe5, 0xa1, 0xab,
	0x12, 0x17, 0x0a, 0x0a, 0x47, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x4d, 0x61, 0x6c, 0x65, 0x10, 0x02,
	0x1a, 0x07, 0x92, 0x9d, 0x20, 0x03, 0xe7, 0x94, 0xb7, 0x12, 0x19, 0x0a, 0x0c, 0x47, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x46, 0x65, 0x6d, 0x61, 0x6c, 0x65, 0x10, 0x03, 0x1a, 0x07, 0x92, 0x9d, 0x20,
	0x03, 0xe5, 0xa5, 0xb3, 0x1a, 0x11, 0xd2, 0xb5, 0x03, 0x0d, 0xf2, 0x01, 0x0a, 0x4f, 0x72, 0x69,
	0x67, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x2a, 0xa3, 0x01, 0x0a, 0x04, 0x52, 0x6f, 0x6c, 0x65,
	0x12, 0x1f, 0x0a, 0x0f, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x68, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x52,
	0x6f, 0x6c, 0x65, 0x10, 0x00, 0x1a, 0x0a, 0x92, 0x9d, 0x20, 0x06, 0xe5, 0x8d, 0xa0, 0xe4, 0xbd,
	0x8d, 0x12, 0x20, 0x0a, 0x0a, 0x52, 0x6f, 0x6c, 0x65, 0x4e, 0x6f, 0x72, 0x6d, 0x61, 0x6c, 0x10,
	0x01, 0x1a, 0x10, 0x92, 0x9d, 0x20, 0x0c, 0xe6, 0x99, 0xae, 0xe9, 0x80, 0x9a, 0xe7, 0x94, 0xa8,
	0xe6, 0x88, 0xb7, 0x12, 0x1c, 0x0a, 0x09, 0x52, 0x6f, 0x6c, 0x65, 0x41, 0x64, 0x6d, 0x69, 0x6e,
	0x10, 0x02, 0x1a, 0x0d, 0x92, 0x9d, 0x20, 0x09, 0xe7, 0xae, 0xa1, 0xe7, 0x90, 0x86, 0xe5, 0x91,
	0x98, 0x12, 0x27, 0x0a, 0x0e, 0x52, 0x6f, 0x6c, 0x65, 0x53, 0x75, 0x70, 0x65, 0x72, 0x41, 0x64,
	0x6d, 0x69, 0x6e, 0x10, 0x03, 0x1a, 0x13, 0x92, 0x9d, 0x20, 0x0f, 0xe8, 0xb6, 0x85, 0xe7, 0xba,
	0xa7, 0xe7, 0xae, 0xa1, 0xe7, 0x90, 0x86, 0xe5, 0x91, 0x98, 0x1a, 0x11, 0xd2, 0xb5, 0x03, 0x0d,
	0xf2, 0x01, 0x0a, 0x4f, 0x72, 0x69, 0x67, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x2a, 0xe4, 0x01,
	0x0a, 0x0a, 0x55, 0x73, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x25, 0x0a, 0x15,
	0x55, 0x73, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x68,
	0x6f, 0x6c, 0x64, 0x65, 0x72, 0x10, 0x00, 0x1a, 0x0a, 0x92, 0x9d, 0x20, 0x06, 0xe5, 0x8d, 0xa0,
	0xe4, 0xbd, 0x8d, 0x12, 0x25, 0x0a, 0x12, 0x55, 0x73, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x49, 0x6e, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x10, 0x01, 0x1a, 0x0d, 0x92, 0x9d, 0x20,
	0x09, 0xe6, 0x9c, 0xaa, 0xe6, 0xbf, 0x80, 0xe6, 0xb4, 0xbb, 0x12, 0x26, 0x0a, 0x13, 0x55, 0x73,
	0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x41, 0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65,
	0x64, 0x10, 0x02, 0x1a, 0x0d, 0x92, 0x9d, 0x20, 0x09, 0xe5, 0xb7, 0xb2, 0xe6, 0xbf, 0x80, 0xe6,
	0xb4, 0xbb, 0x12, 0x23, 0x0a, 0x10, 0x55, 0x73, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x46, 0x72, 0x6f, 0x7a, 0x65, 0x6e, 0x10, 0x03, 0x1a, 0x0d, 0x92, 0x9d, 0x20, 0x09, 0xe5, 0xb7,
	0xb2, 0xe5, 0x86, 0xbb, 0xe7, 0xbb, 0x93, 0x12, 0x24, 0x0a, 0x11, 0x55, 0x73, 0x65, 0x72, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x10, 0x04, 0x1a, 0x0d,
	0x92, 0x9d, 0x20, 0x09, 0xe5, 0xb7, 0xb2, 0xe6, 0xb3, 0xa8, 0xe9, 0x94, 0x80, 0x1a, 0x15, 0xd2,
	0xb5, 0x03, 0x0d, 0xf2, 0x01, 0x0a, 0x4f, 0x72, 0x69, 0x67, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67,
	0xe0, 0xa4, 0x1e, 0x00, 0x2a, 0xc7, 0x02, 0x0a, 0x07, 0x55, 0x73, 0x65, 0x72, 0x45, 0x72, 0x72,
	0x12, 0x22, 0x0a, 0x12, 0x55, 0x73, 0x65, 0x72, 0x45, 0x72, 0x72, 0x50, 0x6c, 0x61, 0x63, 0x65,
	0x68, 0x6f, 0x6c, 0x64, 0x65, 0x72, 0x10, 0x00, 0x1a, 0x0a, 0x92, 0x9d, 0x20, 0x06, 0xe5, 0x8d,
	0xa0, 0xe4, 0xbd, 0x8d, 0x12, 0x2f, 0x0a, 0x0c, 0x55, 0x73, 0x65, 0x72, 0x45, 0x72, 0x72, 0x4c,
	0x6f, 0x67, 0x69, 0x6e, 0x10, 0xe8, 0x07, 0x1a, 0x1c, 0x92, 0x9d, 0x20, 0x18, 0xe7, 0x94, 0xa8,
	0xe6, 0x88, 0xb7, 0xe5, 0x90, 0x8d, 0xe6, 0x88, 0x96, 0xe5, 0xaf, 0x86, 0xe7, 0xa0, 0x81, 0xe9,
	0x94, 0x99, 0xe8, 0xaf, 0xaf, 0x12, 0x29, 0x0a, 0x0f, 0x55, 0x73, 0x65, 0x72, 0x45, 0x72, 0x72,
	0x4e, 0x6f, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x10, 0xe9, 0x07, 0x1a, 0x13, 0x92, 0x9d, 0x20,
	0x0f, 0xe6, 0x9c, 0xaa, 0xe6, 0xbf, 0x80, 0xe6, 0xb4, 0xbb, 0xe8, 0xb4, 0xa6, 0xe5, 0x8f, 0xb7,
	0x12, 0x26, 0x0a, 0x12, 0x55, 0x73, 0x65, 0x72, 0x45, 0x72, 0x72, 0x4e, 0x6f, 0x41, 0x75, 0x74,
	0x68, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x10, 0xea, 0x07, 0x1a, 0x0d, 0x92, 0x9d, 0x20, 0x09, 0xe6,
	0x97, 0xa0, 0xe6, 0x9d, 0x83, 0xe9, 0x99, 0x90, 0x12, 0x2a, 0x0a, 0x13, 0x55, 0x73, 0x65, 0x72,
	0x45, 0x72, 0x72, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x10,
	0xeb, 0x07, 0x1a, 0x10, 0x92, 0x9d, 0x20, 0x0c, 0xe7, 0x99, 0xbb, 0xe5, 0xbd, 0x95, 0xe8, 0xb6,
	0x85, 0xe6, 0x97, 0xb6, 0x12, 0x29, 0x0a, 0x13, 0x55, 0x73, 0x65, 0x72, 0x45, 0x72, 0x72, 0x49,
	0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x10, 0xec, 0x07, 0x1a, 0x0f,
	0x92, 0x9d, 0x20, 0x0b, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0xe9, 0x94, 0x99, 0xe8, 0xaf, 0xaf, 0x12,
	0x22, 0x0a, 0x0e, 0x55, 0x73, 0x65, 0x72, 0x45, 0x72, 0x72, 0x4e, 0x6f, 0x4c, 0x6f, 0x67, 0x69,
	0x6e, 0x10, 0xed, 0x07, 0x1a, 0x0d, 0x92, 0x9d, 0x20, 0x09, 0xe6, 0x9c, 0xaa, 0xe7, 0x99, 0xbb,
	0xe5, 0xbd, 0x95, 0x1a, 0x19, 0xd2, 0xb5, 0x03, 0x0d, 0xf2, 0x01, 0x0a, 0x4f, 0x72, 0x69, 0x67,
	0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0xe0, 0xa4, 0x1e, 0x00, 0xe8, 0xa4, 0x1e, 0x01, 0x32, 0x9b,
	0x01, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x8b,
	0x01, 0x0a, 0x06, 0x53, 0x69, 0x67, 0x6e, 0x75, 0x70, 0x12, 0x0f, 0x2e, 0x75, 0x73, 0x65, 0x72,
	0x2e, 0x53, 0x69, 0x67, 0x6e, 0x75, 0x70, 0x52, 0x65, 0x71, 0x1a, 0x1c, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x52, 0x92, 0x41, 0x32, 0x0a, 0x12, 0xe7,
	0x94, 0xa8, 0xe6, 0x88, 0xb7, 0xe7, 0x9b, 0xb8, 0xe5, 0x85, 0xb3, 0xe6, 0x8e, 0xa5, 0xe5, 0x8f,
	0xa3, 0x0a, 0x06, 0x76, 0x31, 0x2e, 0x30, 0x2e, 0x30, 0x12, 0x06, 0xe6, 0xb3, 0xa8, 0xe5, 0x86,
	0x8c, 0x1a, 0x0c, 0xe6, 0xb3, 0xa8, 0xe5, 0x86, 0x8c, 0xe6, 0x8e, 0xa5, 0xe5, 0x8f, 0xa3, 0xb2,
	0xe0, 0x1f, 0x02, 0x08, 0x01, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x11, 0x3a, 0x01, 0x2a, 0x22, 0x0c,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x42, 0x61, 0xc8, 0x3e,
	0x01, 0x92, 0x41, 0x07, 0x12, 0x05, 0x32, 0x03, 0x31, 0x2e, 0x30, 0xd2, 0xb5, 0x03, 0x02, 0x50,
	0x01, 0x0a, 0x1b, 0x63, 0x68, 0x65, 0x72, 0x72, 0x79, 0x2e, 0x68, 0x6f, 0x70, 0x65, 0x69, 0x6f,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x5a, 0x2f,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x6f, 0x70, 0x65, 0x69,
	0x6f, 0x2f, 0x63, 0x68, 0x65, 0x72, 0x72, 0x79, 0x2f, 0x5f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_user_user_proto_rawDescOnce sync.Once
	file_user_user_proto_rawDescData = file_user_user_proto_rawDesc
)

func file_user_user_proto_rawDescGZIP() []byte {
	file_user_user_proto_rawDescOnce.Do(func() {
		file_user_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_user_user_proto_rawDescData)
	})
	return file_user_user_proto_rawDescData
}

var file_user_user_proto_enumTypes = make([]protoimpl.EnumInfo, 4)
var file_user_user_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_user_user_proto_goTypes = []interface{}{
	(Gender)(0),                    // 0: user.Gender
	(Role)(0),                      // 1: user.Role
	(UserStatus)(0),                // 2: user.UserStatus
	(UserErr)(0),                   // 3: user.UserErr
	(*User)(nil),                   // 4: user.User
	(*SignupReq)(nil),              // 5: user.SignupReq
	(*wrapperspb.StringValue)(nil), // 6: google.protobuf.StringValue
}
var file_user_user_proto_depIdxs = []int32{
	0, // 0: user.User.gender:type_name -> user.Gender
	1, // 1: user.User.role:type_name -> user.Role
	2, // 2: user.User.status:type_name -> user.UserStatus
	0, // 3: user.SignupReq.gender:type_name -> user.Gender
	5, // 4: user.UserService.Signup:input_type -> user.SignupReq
	6, // 5: user.UserService.Signup:output_type -> google.protobuf.StringValue
	5, // [5:6] is the sub-list for method output_type
	4, // [4:5] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_user_user_proto_init() }
func file_user_user_proto_init() {
	if File_user_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_user_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
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
		file_user_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignupReq); i {
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
			RawDescriptor: file_user_user_proto_rawDesc,
			NumEnums:      4,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_user_user_proto_goTypes,
		DependencyIndexes: file_user_user_proto_depIdxs,
		EnumInfos:         file_user_user_proto_enumTypes,
		MessageInfos:      file_user_user_proto_msgTypes,
	}.Build()
	File_user_user_proto = out.File
	file_user_user_proto_rawDesc = nil
	file_user_user_proto_goTypes = nil
	file_user_user_proto_depIdxs = nil
}
