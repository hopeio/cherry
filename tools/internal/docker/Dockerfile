ARG IMAGE=golang:latest

FROM ${IMAGE}

ENV PATH /protoc/bin:$PATH
ENV GOFLAGS -buildvcs=false

RUN git clone https://github.com/hopeio/cherry
RUN cd cherry
RUN go mod download
RUN go mod tidy
RUN cd tools/protoc
RUN echo "开始安装"
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go
RUN protoc -I$protoDir --go_out=paths=source_relative:$cherry/.. $protoDir/cherry/protobuf/utils/**/*.proto
RUN go install ./protoc-gen-enum
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
RUN go install github.com/alta/protopatch/cmd/protoc-gen-go-patch
RUN go install ./protoc-gen-grpc-gin
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
RUN go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators
RUN go install ./protoc-gen-go-patch
#RUN  go install github.com/danielvladco/go-proto-gql/protoc-gen-gql
#RUN  go install github.com/danielvladco/go-proto-gql/protoc-gen-gogql
RUN go install ./protoc-gen-gql
RUN go install ./protoc-gen-gogql
RUN go install github.com/99designs/gqlgen
RUN go install ./protogen
RUN echo "安装完成"
ADD protobuf/_proto/ /proto/
RUN ./install_protoc.sh
RUN mv /protoc


