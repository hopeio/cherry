ARG IMAGE=golang:latest

FROM ${IMAGE}

ENV PATH /protoc/bin:$PATH
ENV GOFLAGS -buildvcs=false


RUN cd / && git clone https://github.com/hopeio/cherry
RUN cd /cherry && go mod download && go mod tidy && copy protobuf/_proto/ /proto/
RUN cd tools/protoc  && ./install_protoc.sh && mv protoc-26.1-linux-x86_64 /protoc && rm protoc-26.1-linux-x86_64.zip
RUN go run /cherry/tools/protoc/install_tools.go
RUN rm /cherry

