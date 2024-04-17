package grpc

import (
	"crypto/tls"
	"github.com/hopeio/cherry/utils/errors/multierr"
	httpi "github.com/hopeio/cherry/utils/net/http"
	"github.com/hopeio/cherry/utils/net/http/grpc/stats"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"strings"
)

var Internal = &metadata.MD{httpi.HeaderInternal: []string{"true"}}

var ClientConns = make(clientConns)

func ClientConnsClose() error {
	if ClientConns != nil {
		return ClientConns.Close()
	}
	return nil
}

type clientConns map[string]*grpc.ClientConn

func (cs clientConns) Close() error {
	var multiErr multierr.MultiError
	for _, conn := range cs {
		err := conn.Close()
		if err != nil {
			multiErr.Append(err)
		}
	}
	if multiErr.HasErrors() {
		return &multiErr
	}
	return nil
}

func GetDefaultClient(addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	if conn, ok := ClientConns[addr]; ok {
		return conn, nil
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithStatsHandler(&stats.InternalClientHandler{}))...)
	if err != nil {
		return nil, err
	}

	ClientConns[addr] = conn
	return conn, nil
}

func GetTlsClient(addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(addr, append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{ServerName: strings.Split(addr, ":")[0], InsecureSkipVerify: true})), grpc.WithStatsHandler(&stats.InternalClientHandler{}))...)
	if err != nil {
		return nil, err
	}
	if oldConn, ok := ClientConns[addr]; ok {
		oldConn.Close()
	}
	ClientConns[addr] = conn
	return conn, nil
}
