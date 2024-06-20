package generator

import (
	gqlpb "github.com/danielvladco/go-proto-gql/pkg/graphqlpb"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/proto"
	descriptor "google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func GraphqlMethodOptions(opts proto.Message) *gqlpb.Rpc {
	if opts != nil {
		v := proto.GetExtension(opts, gqlpb.E_Rpc)
		if v != nil {
			return v.(*gqlpb.Rpc)
		}
	}
	return nil
}

func GraphqlServiceOptions(opts proto.Message) *gqlpb.Svc {
	if opts != nil {
		v := proto.GetExtension(opts, gqlpb.E_Svc)
		if v != nil {
			return v.(*gqlpb.Svc)
		}
	}
	return nil
}

func GraphqlFieldOptions(opts proto.Message) *gqlpb.Field {
	if opts != nil {
		v := proto.GetExtension(opts, gqlpb.E_Field)
		if v != nil && v.(*gqlpb.Field) != nil {
			return v.(*gqlpb.Field)
		}
	}
	return nil
}

func GraphqlOneofOptions(opts proto.Message) *gqlpb.Oneof {
	if opts != nil {
		v := proto.GetExtension(opts, gqlpb.E_Oneof)
		if v != nil && v.(*gqlpb.Oneof) != nil {
			return v.(*gqlpb.Oneof)
		}
	}
	return nil
}

func GetRequestType(rpcOpts *gqlpb.Rpc, svcOpts *gqlpb.Svc) gqlpb.Type {
	if rpcOpts != nil && rpcOpts.Type != nil {
		return *rpcOpts.Type
	}
	if svcOpts != nil && svcOpts.Type != nil {
		return *svcOpts.Type
	}
	return gqlpb.Type_DEFAULT
}

func CreateDescriptorsFromProto(req *pluginpb.CodeGeneratorRequest) (descs []*desc.FileDescriptor, err error) {
	dd, err := desc.CreateFileDescriptorsFromSet(&descriptor.FileDescriptorSet{
		File: req.GetProtoFile(),
	})
	if err != nil {
		return nil, err
	}
	for _, d := range dd {
		for _, filename := range req.FileToGenerate {
			if filename == d.GetName() {
				descs = append(descs, d)
			}
		}
	}
	return
}

func ResolveProtoFilesRecursively(descs []*desc.FileDescriptor) (files FileDescriptors) {
	for _, f := range descs {
		files = append(files, ResolveProtoFilesRecursively(f.GetDependencies())...)
		files = append(files, f)
	}

	return files
}

type FileDescriptors []*desc.FileDescriptor

func (ds FileDescriptors) AsFileDescriptorProto() (files []*descriptor.FileDescriptorProto) {
	for _, d := range ds {
		files = append(files, d.AsFileDescriptorProto())
	}
	return
}

func baseType(field *desc.FieldDescriptor) bool {
	switch field.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE,
		descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return true

	case descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_SINT64,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_SINT32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_FIXED32,
		descriptor.FieldDescriptorProto_TYPE_UINT64,
		descriptor.FieldDescriptorProto_TYPE_FIXED64:
		return true

	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return true

	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return true
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		return true
	default:
		return false
	}
}
