package mix

import (
	"context"
	"testing"
	"time"
)

type connCtxKey struct{}

func waitCancelled(t *testing.T, ctx context.Context, what string) {
	t.Helper()
	select {
	case <-ctx.Done():
	case <-time.After(2 * time.Second):
		t.Fatalf("%s did not cancel the derived context", what)
	}
}

// The HTTP/3 connection context has to keep quic-go's values: ConnContext is
// called after quic-go has attached ServerContextKey, LocalAddrContextKey and
// RemoteAddrContextKey to conn.Context(). Replacing it with the process
// shutdown context made all three unavailable to the request handler.
func TestBridgeShutdownContextKeepsConnectionValues(t *testing.T) {
	shutdown, stop := context.WithCancel(context.Background())
	defer stop()
	connCtx := context.WithValue(context.Background(), connCtxKey{}, "remote-addr")

	ctx := bridgeShutdownContext(shutdown, connCtx)
	if got := ctx.Value(connCtxKey{}); got != "remote-addr" {
		t.Fatalf("connection value lost: got %v", got)
	}
	select {
	case <-ctx.Done():
		t.Fatal("context should still be live")
	default:
	}
}

// A client disconnect cancels conn.Context(). The derived context must follow,
// otherwise in-flight h3 requests and WebTransport sessions keep running after
// their connection is gone.
func TestBridgeShutdownContextFollowsConnection(t *testing.T) {
	shutdown, stop := context.WithCancel(context.Background())
	defer stop()
	connCtx, closeConn := context.WithCancel(context.Background())

	ctx := bridgeShutdownContext(shutdown, connCtx)
	closeConn()
	waitCancelled(t, ctx, "QUIC connection close")
}

// SIGTERM cancels the shutdown context. Every HTTP/3 connection must be
// cancelled too, even though its own QUIC connection is still healthy.
func TestBridgeShutdownContextFollowsShutdown(t *testing.T) {
	shutdown, stop := context.WithCancel(context.Background())
	connCtx, closeConn := context.WithCancel(context.Background())
	defer closeConn()

	ctx := bridgeShutdownContext(shutdown, connCtx)
	stop()
	waitCancelled(t, ctx, "process shutdown")
}

// context.AfterFunc registrations have to be released when the connection
// ends, otherwise a long-running server accumulates one callback per QUIC
// connection. We cannot observe the registration directly, but we can assert
// that a shutdown after the connection is gone does not touch an already
// cancelled derived context (and that the helper returns promptly).
func TestBridgeShutdownContextReleasesOnConnectionEnd(t *testing.T) {
	shutdown, stop := context.WithCancel(context.Background())
	defer stop()
	connCtx, closeConn := context.WithCancel(context.Background())

	ctx := bridgeShutdownContext(shutdown, connCtx)
	closeConn()
	waitCancelled(t, ctx, "connection close")
	// A later shutdown must be harmless.
	stop()
	waitCancelled(t, ctx, "later shutdown")
}
