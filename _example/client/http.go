package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/hopeio/cherry/_example/protobuf/user"
	"golang.org/x/net/http2"
	"google.golang.org/protobuf/proto"
)

func main() {
	sendGRPCRequest()
}

func createGRPCRequest(msg proto.Message) ([]byte, error) {
	// 1. 序列化 protobuf
	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	// 2. 创建 gRPC 帧
	frame := make([]byte, 5+len(data))

	// 3. 压缩标志位 (0 = 未压缩)
	frame[0] = 0

	// 4. 写入消息长度 (大端序)
	binary.BigEndian.PutUint32(frame[1:5], uint32(len(data)))

	// 5. 写入消息数据
	copy(frame[5:], data)

	return frame, nil
}

func sendGRPCRequest() {
	// 准备请求消息
	req := &user.GetUserReq{Id: 1}

	// 创建 gRPC 帧
	grpcBody, _ := createGRPCRequest(req)

	// 创建 HTTP 请求
	httpReq, _ := http.NewRequest(
		"POST",
		"http://localhost:8080/user.UserService/GetUser",
		bytes.NewReader(grpcBody),
	)

	// 设置必需的 gRPC 头部
	httpReq.Header.Set("Content-Type", "application/grpc+proto")
	httpReq.Header.Set("Te", "trailers") // 必须小写 "te"
	httpReq.Header.Set("User-Agent", "grpc-go/1.0")
	httpReq.ProtoMinor = 2
	// 发送请求
	client := &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
				// 忽略 TLS，使用普通 TCP
				return net.Dial(network, addr)
			},
		},
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, _ := io.ReadAll(resp.Body)

	// 解析 gRPC 响应帧
	if len(body) < 5 {
		panic("invalid response")
	}

	// 检查压缩标志
	compressedFlag := body[0]

	// 读取消息长度
	msgLength := binary.BigEndian.Uint32(body[1:5])

	// 提取消息数据
	msgData := body[5 : 5+msgLength]

	// 解析响应
	var response user.User
	if compressedFlag == 1 {
		// 需要解压
		panic("compressed response not supported")
	}
	proto.Unmarshal(msgData, &response)

	fmt.Printf("Response: %v\n", &response)

	// 检查 gRPC 状态（在 trailers 中）
	grpcStatus := resp.Trailer.Get("Grpc-Status")
	if grpcStatus != "" && grpcStatus != "0" {
		grpcMessage := resp.Trailer.Get("Grpc-Message")
		fmt.Printf("gRPC error: %s - %s\n", grpcStatus, grpcMessage)
	}
}
