/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT license that can be found in the LICENSE file.
 * @Created by jyb
 */

package mix

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
)

// Behaviour when mix owns the QUIC stack.
//
// This pins down that WebTransport and plain HTTP/3 coexist on one UDP
// port: the same port serves regular h3 requests and establishes WT
// sessions. There is only one way to do that: webtransport.Server owns
// the QUIC stack. If http3.Server kept its own listener while WT started
// a second one, each connection would end up with two HTTP/3 control
// streams and the client would abort with "duplicate control stream"
// (RFC 9114 section 6.2.1 requires exactly one per connection).

// TestWebTransportSharesPortWithHTTP3 runs regular h3 requests and WT
// sessions on the same port.
func TestWebTransportSharesPortWithHTTP3(t *testing.T) {
	certFile, keyFile := writeTestCert(t)
	addr := freeUDPAddr(t)

	// Session handler: read one message from the session and reply.
	sessCh := make(chan string, 1)
	// Regular h3 route: shares the port with WT without interfering.
	mux := http.NewServeMux()
	mux.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello-h3")
	})
	// The application routes the WebTransport CONNECT itself and upgrades
	// it through mix, the same way any application would.
	mux.HandleFunc("/api/wt", func(w http.ResponseWriter, r *http.Request) {
		sess, err := UpgradeWebTransport(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		go func() {
			stream, err := sess.AcceptStream(context.Background())
			if err != nil {
				return
			}
			buf := make([]byte, 16)
			n, _ := stream.Read(buf)
			sessCh <- string(buf[:n])
			// Reply so the client can verify both directions work.
			_, _ = stream.Write([]byte("pong"))
		}()
	})
	s := &Server{
		Addr:        addr,
		CertFile:    certFile,
		KeyFile:     keyFile,
		HttpHandler: mux,
	}
	s.HTTP3.Enabled = true
	s.HTTP3.Addr = addr
	s.HTTP3.CertFile = certFile
	s.HTTP3.KeyFile = keyFile
	s.HTTP3.WebTransport.Enabled = true

	// mix's Run manages the signal context itself (SIGINT/SIGTERM) and
	// returns no error.
	go s.Run()

	// Wait for the listener to come up.
	waitForListener(t, addr)

	// 1) A regular HTTP/3 request, on the same port as WT.
	if got := h3Get(t, addr); got != "hello-h3" {
		t.Fatalf("h3 response = %q, want %q", got, "hello-h3")
	}

	// 2) A WebTransport session, on the same port.
	client := dialWT(t, addr)
	stream, err := client.OpenStream()
	if err != nil {
		t.Fatalf("open stream: %v", err)
	}
	if _, err := stream.Write([]byte("hello-wt")); err != nil {
		t.Fatalf("write: %v", err)
	}
	select {
	case got := <-sessCh:
		if got != "hello-wt" {
			t.Fatalf("session got %q, want %q", got, "hello-wt")
		}
	case <-time.After(10 * time.Second):
		t.Fatal("WT session handler not called")
	}
	// Read the bytes the server sent back.
	buf := make([]byte, 4)
	if _, err := io.ReadFull(stream, buf); err != nil {
		t.Fatalf("read reply: %v", err)
	}
	if string(buf) != "pong" {
		t.Fatalf("reply = %q, want %q", buf, "pong")
	}
}

// TestIsWebTransportRequest checks that request detection looks at both
// the method and the protocol.
func TestIsWebTransportRequest(t *testing.T) {
	cases := []struct {
		name   string
		method string
		proto  string
		want   bool
	}{
		{"WT 的 CONNECT", http.MethodConnect, "webtransport", true},
		{"普通 CONNECT", http.MethodConnect, "h3", false},
		{"普通 GET", http.MethodGet, "webtransport", false},
		{"普通 GET h3", http.MethodGet, "HTTP/3.0", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := &http.Request{Method: tc.method, Proto: tc.proto}
			if got := isWebTransportRequest(r); got != tc.want {
				t.Fatalf("isWebTransportRequest = %v, want %v", got, tc.want)
			}
		})
	}
}

// h3Get sends a regular HTTP/3 request to addr.
func h3Get(t *testing.T, addr string) string {
	t.Helper()
	client := &http.Client{
		Transport: &http3.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var lastErr error
	for ctx.Err() == nil {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://"+addr+"/api/ping", nil)
		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			return string(body)
		}
		lastErr = err
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("h3 request never succeeded: %v", lastErr)
	return ""
}

// dialWT establishes a WebTransport session.
func dialWT(t *testing.T, addr string) *webtransport.Session {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var lastErr error
	for ctx.Err() == nil {
		qconn, err := quic.DialAddr(ctx, addr,
			&tls.Config{InsecureSkipVerify: true, NextProtos: []string{"h3"}},
			&quic.Config{EnableDatagrams: true, EnableStreamResetPartialDelivery: true})
		if err != nil {
			lastErr = err
			time.Sleep(50 * time.Millisecond)
			continue
		}
		tr := &webtransport.Transport{}
		cc, err := tr.NewClientConn(qconn)
		if err != nil {
			lastErr = err
			time.Sleep(50 * time.Millisecond)
			continue
		}
		dialCtx, dialCancel := context.WithTimeout(ctx, 3*time.Second)
		_, sess, err := cc.Dial(dialCtx, "https://"+addr+"/api/wt", nil)
		dialCancel()
		if err != nil {
			lastErr = err
			time.Sleep(50 * time.Millisecond)
			continue
		}
		return sess
	}
	t.Fatalf("webtransport session never established: %v", lastErr)
	return nil
}

// waitForListener waits until the UDP port answers.
func waitForListener(t *testing.T, addr string) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.Dial("udp", addr)
		if err == nil {
			_ = conn.Close()
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("listener on %s never came up", addr)
}

// freeUDPAddr reserves a free local UDP port.
func freeUDPAddr(t *testing.T) string {
	t.Helper()
	conn, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve udp port: %v", err)
	}
	addr := conn.LocalAddr().String()
	_ = conn.Close()
	return addr
}

// writeTestCert generates a self-signed certificate.
func writeTestCert(t *testing.T) (certFile, keyFile string) {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	tpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "localhost"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		DNSNames:              []string{"localhost"},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IsCA:                  true,
		BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, tpl, tpl, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("create certificate: %v", err)
	}
	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		t.Fatalf("marshal key: %v", err)
	}
	dir := t.TempDir()
	certFile = filepath.Join(dir, "cert.pem")
	keyFile = filepath.Join(dir, "key.pem")
	_ = os.WriteFile(certFile, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0o600)
	_ = os.WriteFile(keyFile, pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}), 0o600)
	return certFile, keyFile
}
