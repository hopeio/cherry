cherry=$(go list -m -f {{.Dir}}  github.com/hopeio/cherry)
cherry=${cherry//\\/\/}
# 复制proto文件
echo "copy proto dependencies ..."
## 依赖地址
gatewayDir=$(go list -m -f {{.Dir}} github.com/grpc-ecosystem/grpc-gateway/v2)
gatewayDir=${gatewayDir//\\/\/}
validatorsDir=$(go list -m -f {{.Dir}} github.com/mwitkow/go-proto-validators)
validatorsDir=${validatorsDir//\\/\/}
protopatchDir=$(go list -m -f {{.Dir}} github.com/alta/protopatch)
protopatchDir=${protopatchDir//\\/\/}
gqlDir=$(go list -m -f {{.Dir}} github.com/danielvladco/go-proto-gql)
gqlDir=${gqlDir//\\/\/}

protoDir=$cherry/protobuf/_proto

## googleapis
cd $protoDir
go mod init proto
go get github.com/googleapis/googleapis
googleapisDir=$(go list -m -f {{.Dir}} github.com/googleapis/googleapis)
googleapisDir=${googleapisDir//\\/\/}
echo $googleapisDir
go get github.com/protocolbuffers/protobuf
protobufDir=$(go list -m -f {{.Dir}} github.com/protocolbuffers/protobuf)
protobufDir=${protobufDir//\\/\/}
echo $protobufDir

## copy
cp  $gatewayDir/protoc-gen-openapiv2/options/*.proto $protoDir/protoc-gen-openapiv2/options
cp  $googleapisDir/google/api/*.proto $protoDir/google/api
cp  $validatorsDir/*.proto $protoDir/github.com/mwitkow/go-proto-validators
cp  $protobufDir/src/google/protobuf/*.proto $protoDir/google/protobuf
不使用github.com/alta/protopatch
cp  $protopatchDir/patch/*.proto $protoDir/patch
cp  $gqlDir/api/danielvladco/protobuf/*.proto $protoDir/danielvladco/protobuf

rm $protoDir/go.mod
rm $protoDir/go.sum
