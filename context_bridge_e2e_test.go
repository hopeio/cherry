package mix

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/quic-go/quic-go/http3"
)

// TestH3RequestKeepsQUICConnectionContext is the end-to-end counterpart of
// the bridgeShutdownContext unit tests. A real HTTP/3 request must arrive at
// the handler carrying the values quic-go attaches to conn.Context().
//
// The old ConnContext returned the process shutdown context instead, so all
// three of these were nil. That is not a cosmetic loss: middleware that reads
// the peer address from the context (rather than r.RemoteAddr, which HTTP/3
// leaves empty) silently saw nothing.
func TestH3RequestKeepsQUICConnectionContext(t *testing.T) {
	certFile, keyFile := writeTestCert(t)
	addr := freeUDPAddr(t)

	type observed struct {
		server  bool
		local   bool
		remote  bool
		liveCtx bool
	}
	got := make(chan observed, 1)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/ctx", func(w http.ResponseWriter, r *http.Request) {
		_, localOK := r.Context().Value(http.LocalAddrContextKey).(net.Addr)
		_, remoteOK := r.Context().Value(http3.RemoteAddrContextKey).(net.Addr)
		got <- observed{
			server:  r.Context().Value(http3.ServerContextKey) != nil,
			local:   localOK,
			remote:  remoteOK,
			liveCtx: r.Context().Err() == nil,
		}
		fmt.Fprint(w, "ok")
	})

	s := &Server{
		Addr:                  addr,
		CertFile:              certFile,
		KeyFile:               keyFile,
		HttpHandler:           mux,
		DisableInternalServer: true,
	}
	s.HTTP3.Enabled = true
	s.HTTP3.Addr = addr
	s.HTTP3.CertFile = certFile
	s.HTTP3.KeyFile = keyFile
	go s.Run()
	waitForListener(t, addr)

	if body := h3GetPath(t, addr, "/api/ctx"); body != "ok" {
		t.Fatalf("h3 response = %q, want ok", body)
	}
	select {
	case o := <-got:
		if !o.server {
			t.Error("http3.ServerContextKey missing from request context")
		}
		if !o.local {
			t.Error("http.LocalAddrContextKey missing from request context")
		}
		if !o.remote {
			t.Error("http3.RemoteAddrContextKey missing from request context")
		}
		if !o.liveCtx {
			t.Error("request context was already cancelled on arrival")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("h3 handler was not called")
	}
}

// h3GetPath is h3Get against a caller-supplied path.
func h3GetPath(t *testing.T, addr, path string) string {
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
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://"+addr+path, nil)
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
