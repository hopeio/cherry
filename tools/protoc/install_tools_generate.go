//go:build tools

package main

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/alta/protopatch/cmd/protoc-gen-go-patch"
	_ "github.com/bufbuild/protovalidate-go"
	_ "github.com/danielvladco/go-proto-gql/pkg/graphqlpb"
	_ "github.com/envoyproxy/protoc-gen-validate"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "github.com/mwitkow/go-proto-validators/protoc-gen-govalidators"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)

// Deprecated
//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go
////go:generate protoc -I../../protobuf/_proto --go_out=paths=source_relative:../../.. ../../protobuf/_proto/cherry/protobuf/utils/patch/*.proto
//go:generate protoc -I../../protobuf/_proto --go_out=paths=source_relative:../../.. ../../protobuf/_proto/cherry/protobuf/utils/apiconfig/*.proto
//go:generate protoc -I../../protobuf/_proto --go_out=paths=source_relative:../../.. ../../protobuf/_proto/cherry/protobuf/utils/openapiconfig/*.proto
//go:generate protoc -I../../protobuf/_proto --go_out=paths=source_relative:../../.. ../../protobuf/_proto/cherry/protobuf/utils/enum/*.proto
//go:generate go install ./protoc-gen-enum
//go:generate go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
//go:generate go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
//go:generate go install github.com/alta/protopatch/cmd/protoc-gen-go-patch
//go:generate go install ./protoc-gen-grpc-gin
//go:generate go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
//go:generate go install github.com/envoyproxy/protoc-gen-validate
////go:generate go install ./protoc-gen-go-patch
////go:generate go install github.com/danielvladco/go-proto-gql/cmd/proto2graphql
////go:generate go install github.com/danielvladco/go-proto-gql/protoc-gen-gql
////go:generate go install github.com/danielvladco/go-proto-gql/protoc-gen-gogql
//go:generate go install ./protoc-gen-gql
//go:generate go install ./protoc-gen-gogql
//go:generate go install github.com/99designs/gqlgen
//go:generate go install ./protogen
//go:generate echo "Installation Finished"
