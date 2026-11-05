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

	"github.com/hopeio/mix/_example/protobuf/user"
	"google.golang.org/protobuf/proto"
)

func main() {
	sendGRPCRequest()
}

func createGRPCRequest(msg proto.Message) ([]byte, error) {
	// 1. Serialize the protobuf message.
	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	// 2. Build the gRPC frame.
	frame := make([]byte, 5+len(data))

	// 3. The compressed flag (0 means not compressed).
	frame[0] = 0

	// 4. Write the message length (big endian).
	binary.BigEndian.PutUint32(frame[1:5], uint32(len(data)))

	// 5. Write the message body.
	copy(frame[5:], data)

	return frame, nil
}

func sendGRPCRequest() {
	// Prepare the request message.
	req := &user.GetUserReq{Id: 1}

	// Build the gRPC frame.
	grpcBody, _ := createGRPCRequest(req)

	// Build the HTTP request.
	httpReq, _ := http.NewRequest(
		"POST",
		"http://localhost:8080/user.UserService/GetUser",
		bytes.NewReader(grpcBody),
	)

	// Set the gRPC headers that are required.
	httpReq.Header.Set("Content-Type", "application/grpc+proto")
	httpReq.Header.Set("Te", "trailers") // The "te" header must be lowercase.
	httpReq.Header.Set("User-Agent", "grpc-go/1.0")
	httpReq.ProtoMinor = 2
	tr := &http.Transport{}
	tr.Protocols = new(http.Protocols)
	tr.Protocols.SetUnencryptedHTTP2(true)
	// Send the request.
	client := &http.Client{
		Transport: tr,
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read the response.
	body, _ := io.ReadAll(resp.Body)

	// Parse the gRPC response frame.
	if len(body) < 5 {
		panic("invalid response")
	}

	// Check the compressed flag.
	compressedFlag := body[0]

	// Read the message length.
	msgLength := binary.BigEndian.Uint32(body[1:5])

	// Extract the message body.
	msgData := body[5 : 5+msgLength]

	// Parse the response.
	var response user.User
	if compressedFlag == 1 {
		// The body has to be decompressed.
		panic("compressed response not supported")
	}
	proto.Unmarshal(msgData, &response)

	fmt.Printf("Response: %v\n", &response)

	// Check stream status in HTTP trailers (Error-Code; same as unary HTTP).
	errCode := resp.Trailer.Get("Error-Code")
	if errCode != "" && errCode != "0" {
		fmt.Printf("HTTP stream error: Error-Code=%s\n", errCode)
	}
}
