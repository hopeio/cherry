cd $(dirname $0) && pwd
cherry=$(go list -m -f {{.Dir}}  github.com/hopeio/cherry)
cherry=${cherry//\\/\/}
protoDir=$cherry/protobuf/_proto

if [ -n "$1" ]; then
  export PATH=$1/bin:$PATH
  echo $PATH
fi

if ! command -v protoc &> /dev/null
then
    echo "protoc command not found,will download"
    # 在这里执行其他操作
    source ./install_protoc.sh
fi


# 安装
cd $cherry/tools/protoc
echo "Start Installation"
go install google.golang.org/protobuf/cmd/protoc-gen-go
protoc -I$protoDir --go_out=paths=source_relative:$cherry/.. $protoDir/cherry/protobuf/utils/**/*.proto
go install ./protoc-gen-enum
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
go install github.com/alta/protopatch/cmd/protoc-gen-go-patch
go install ./protoc-gen-grpc-gin
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators
go install ./protoc-gen-go-patch
# go install github.com/danielvladco/go-proto-gql/protoc-gen-gql
# go install github.com/danielvladco/go-proto-gql/protoc-gen-gogql
go install ./protoc-gen-gql
go install ./protoc-gen-gogql
go install github.com/99designs/gqlgen
go install ./protogen
echo "Installation Finished"
