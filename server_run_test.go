package mix

import (
	"context"
	"io"
	"net"
	"net/http"
	"testing"
	"time"
)

// Run used to log.Fatalf on a listener failure, which skipped its own shutdown pass and
// every defer the caller had registered. These tests pin what the error return buys:
// the failure comes back to the caller instead of ending the process.

func TestRunReturnsListenerErrorInsteadOfFataling(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	// Held for the whole test, so Run's own bind must fail with EADDRINUSE.
	defer func() { _ = ln.Close() }()

	s := &Server{}
	s.Addr = addr
	s.DisableInternalServer = true

	done := make(chan error, 1)
	go func() { done <- s.Run() }()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("a listener that cannot bind must hand back an error, not stop silently")
		}
	case <-time.After(20 * time.Second):
		t.Fatal("Run did not return after the fatal listener error")
	}

	// This test is the point: under the old code the Fatal in this path killed the test
	// binary outright (no FAIL line, just a vanished process), which is the same way it
	// skipped the caller's defers. Now Run returns, and its own shutdown pass has
	// already completed — a Shutdown call once that state is reached is an idempotent nil.
	if err := s.Shutdown(context.Background()); err != nil {
		t.Fatalf("expected the shutdown pass to have completed, got %v", err)
	}
}

// An optional transport that cannot start must not take the API down. HTTP/3 is given a
// certificate pair that cannot load, which is exactly the failure prepareHTTP3 does not
// catch (the fields are set, the material is broken).
func TestRunSurvivesOptionalHTTP3Failure(t *testing.T) {
	addr := freeAddr(t)

	s := &Server{}
	s.Addr = addr
	s.DisableInternalServer = true
	s.HTTP3.Enabled = true
	s.HTTP3.Addr = "127.0.0.1:0"
	s.HTTP3.CertFile, s.HTTP3.KeyFile = "missing.crt", "missing.key"

	ctx, cancel := context.WithCancel(context.Background())
	s.WithContext(ctx)

	done := make(chan error, 1)
	go func() { done <- s.Run() }()

	// The main listener must actually serve while the optional one is failing.
	if !waitServing(t, addr) {
		cancel()
		t.Fatalf("main listener never served on %s", addr)
	}
	cancel() // cancelling the base context is the same path a SIGTERM takes

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("an optional transport failing must not end Run with an error: %v", err)
		}
	case <-time.After(30 * time.Second):
		t.Fatal("Run did not return after the shutdown signal")
	}
}

// freeAddr returns an address that was bindable a moment ago. Used for the main
// listener, which must succeed; the racing window is acceptable in a test that fails
// loudly if it loses.
func freeAddr(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	_ = ln.Close()
	return addr
}

func waitServing(t *testing.T, addr string) bool {
	t.Helper()
	// The target is the 127.0.0.1 port this test just bound; no external input involved.
	for i := 0; i < 100; i++ {
		resp, err := http.Get("http://" + addr + "/")
		if err == nil {
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			return true
		}
		time.Sleep(50 * time.Millisecond)
	}
	return false
}
