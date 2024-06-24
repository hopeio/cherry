cherry=$(go list -m -f {{.Dir}}  github.com/hopeio/cherry)
cherry=${cherry//\\/\/}
echo $cherry
# 复制proto文件
echo "copy proto dependencies ..."
## 依赖地址
gatewayDir=$(go list -m -f {{.Dir}} github.com/grpc-ecosystem/grpc-gateway/v2)
gatewayDir=${gatewayDir//\\/\/}
#validatorsDir=$(go list -m -f {{.Dir}} github.com/envoyproxy/protoc-gen-validate)
#validatorsDir=${validatorsDir//\\/\/}

protopatchDir=$(go list -m -f {{.Dir}} github.com/alta/protopatch)
protopatchDir=${protopatchDir//\\/\/}
gqlDir=$(go list -m -f {{.Dir}} github.com/danielvladco/go-proto-gql)
gqlDir=${gqlDir//\\/\/}

protoDir=$cherry/protobuf/_proto

tmpMod=/tmp/proto
mkdir $tmpMod
## googleapis
cd $tmpMod
go mod init proto
go mod tidy
go get github.com/googleapis/googleapis
googleapisDir=$(go list -m -f {{.Dir}} github.com/googleapis/googleapis)
googleapisDir=${googleapisDir//\\/\/}
echo $googleapisDir
go get github.com/protocolbuffers/protobuf
protobufDir=$(go list -m -f {{.Dir}} github.com/protocolbuffers/protobuf)
protobufDir=${protobufDir//\\/\/}
echo $protobufDir
go get github.com/bufbuild/protovalidate@main
validatorsDir=$(go list -m -f {{.Dir}} github.com/bufbuild/protovalidate)
validatorsDir=${validatorsDir//\\/\/}
echo $validatorsDir
## copy
mkdir -p $protoDir/protoc-gen-openapiv2/options
cp  $gatewayDir/protoc-gen-openapiv2/options/*.proto $protoDir/protoc-gen-openapiv2/options
mkdir -p $protoDir/google/api
cp  $googleapisDir/google/api/*.proto $protoDir/google/api
mkdir -p $protoDir/buf/validate/priv
cp  $validatorsDir/proto/protovalidate/buf/validate/*.proto $protoDir/buf/validate
# chmod -R 777 $validatorsDir/proto/protovalidate/buf/validate/priv
cp  $validatorsDir/proto/protovalidate/buf/validate/priv/*.proto $protoDir/buf/validate/priv/
mkdir -p $protoDir/google/protobuf
cp  $protobufDir/src/google/protobuf/*.proto $protoDir/google/protobuf
rm $protoDir/google/protobuf/unittest*.proto
rm $protoDir/google/protobuf/test_*.proto
rm $protoDir/google/protobuf/*unittest.proto
rm $protoDir/google/protobuf/*_test.proto
rm $protoDir/google/protobuf/sample_messages_edition.proto.proto
#不使用github.com/alta/protopatch
mkdir -p $protoDir/patch
cp  $protopatchDir/patch/*.proto $protoDir/patch
mkdir -p $protoDir/danielvladco/protobuf
cp  $gqlDir/api/danielvladco/protobuf/*.proto $protoDir/danielvladco/protobuf

rm -r $tmpMod