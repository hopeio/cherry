package mix

import (
	"crypto/tls"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
)

// WebTransport support.
//
// WebTransport is an extension of HTTP/3 (RFC 9297): the client opens a
// session with an extended CONNECT (":protocol=webtransport") on the same
// QUIC connection, after which both sides can open bidirectional streams
// inside that session.
//
// Why mix owns the QUIC stack: WebTransport and plain HTTP/3 share one UDP
// port, and webtransport-go has to decode regular h3 requests and WT
// CONNECT requests on the same connection. If http3.Server kept its own
// listener while WT started a second one, each connection would end up with
// two HTTP/3 control streams and the client would abort with "duplicate
// control stream" (RFC 9114 section 6.2.1 requires exactly one per
// connection).
//
// So when WebTransport is enabled, mix hands the QUIC stack to
// webtransport.Server and forwards regular h3 requests to the application
// handler. The WT CONNECT request reaches the application like any other
// request: the application routes it (typically to "/api/wt") and calls
// [UpgradeWebTransport] to turn it into a session.
//
// mix does not decide where WebTransport lives or what happens to a
// session: that is the application's business.

// WebTransportConfig configures WebTransport.
type WebTransportConfig struct {
	// Enabled turns WebTransport on.
	//
	// When off (the default) mix only serves plain HTTP/3 and none of the
	// WebTransport behaviour kicks in. Clients fall back to WebSocket
	// automatically, so the feature set stays complete.
	Enabled bool
	// CheckOrigin validates the cross-site origin. nil disables the
	// check (local debugging only).
	//
	// Browsers send Origin on WT sessions. Without the check, any site
	// could use the user's browser to open a WebSocket-like connection to
	// this service (cross-site hijacking). Non-browser clients send no
	// Origin and are allowed through.
	CheckOrigin func(r *http.Request) bool
	// SessionTimeout is the idle timeout for a session; <=0 uses the
	// default.
	SessionTimeout time.Duration
}

// Default idle timeout for a session: one minute. WebTransport sessions are
// usually long-lived, but an idle one should still be reclaimed.
const defaultWebTransportSessionTimeout = 60 * time.Second

// webTransportProtocols are the accepted :protocol values.
//
// RFC 9297 defines "webtransport"; earlier drafts used "webtransport-h3".
// webtransport-go accepts both, and clients differ (browsers send the
// former, some libraries the latter), so the server has to take both.
var webTransportProtocols = []string{"webtransport", "webtransport-h3"}

// isWebTransportRequest reports whether the request is a WebTransport
// CONNECT.
//
// Two conditions: the method is CONNECT and the protocol is webtransport
// (or its legacy alias). Missing either means it is not a WT request.
func isWebTransportRequest(r *http.Request) bool {
	if r.Method != http.MethodConnect {
		return false
	}
	for _, p := range webTransportProtocols {
		if r.Proto == p {
			return true
		}
	}
	return false
}

// webTransportServer holds the WebTransport server of the running HTTP/3
// listener.
//
// It is package level because the application upgrades sessions from its
// own request handler, which has no reference to the mix Server. A process
// serves one HTTP/3 listener, so there is no ambiguity.
var (
	webTransportMu     sync.RWMutex
	webTransportServer *webtransport.Server
)

// UpgradeWebTransport upgrades a WebTransport CONNECT request into a
// session.
//
// The application calls it from its own route handler, typically:
//
//	app.Handle("/api/wt", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		sess, err := mix.UpgradeWebTransport(w, r)
//		...
//	}))
//
// It fails when WebTransport is not enabled or the request is not a
// WebTransport CONNECT.
//
// On success the connection has been hijacked: the caller owns the session
// from then on and is responsible for closing it.
func UpgradeWebTransport(w http.ResponseWriter, r *http.Request) (*webtransport.Session, error) {
	if !isWebTransportRequest(r) {
		return nil, errors.New("mix: not a webtransport request")
	}
	webTransportMu.RLock()
	wt := webTransportServer
	webTransportMu.RUnlock()
	if wt == nil {
		return nil, errors.New("mix: webtransport is not enabled")
	}
	return wt.Upgrade(w, r)
}

// newWebTransportServer builds a WebTransport server and attaches it to h3.
//
// The important part is ConfigureHTTP3Server: it installs the SETTINGS
// WebTransport needs and injects the QUIC connection into the request
// context. Without it Upgrade cannot reach the underlying connection.
func newWebTransportServer(h3 *http3.Server, cfg WebTransportConfig) (*webtransport.Server, error) {
	wt := &webtransport.Server{
		CheckOrigin: cfg.CheckOrigin,
	}
	webtransport.ConfigureHTTP3Server(h3)
	wt.H3 = h3
	return wt, nil
}

// serveHTTP3 starts HTTP/3, optionally with WebTransport.
//
// Two modes:
//   - WebTransport off: http3.Server listens on its own.
//   - WebTransport on: webtransport.Server owns the QUIC stack, and
//     regular h3 requests and WT sessions coexist on the same connection.
//     It is the only way both can share one port.
//
// The caller only reaches this when [Server.HTTP3].Enabled is set, so
// WebTransport implies HTTP/3: turning on HTTP/3 alone is enough for the
// WebTransport switch to take effect.
//
// The certificate is shared with the HTTPS listener: both point at the same
// cert and key files, so a deployment only has one certificate to manage.
// The tls.Config itself cannot be shared, because the two listeners need
// different ALPN sets (h2/http/1.1 over TCP, h3 over QUIC) and
// http3.ConfigureTLSConfig adds the h3 entry.
func (s *Server) serveHTTP3(certFile, keyFile string, handler http.Handler) error {
	if s.HTTP3.TLSConfig == nil {
		// The certificate was not set up (no cert and key configured).
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return err
		}
		s.HTTP3.TLSConfig = http3.ConfigureTLSConfig(&tls.Config{Certificates: []tls.Certificate{cert}})
	}

	cfg := s.HTTP3.WebTransport
	if !cfg.Enabled {
		return s.HTTP3.ListenAndServeTLS(certFile, keyFile)
	}

	wt, err := newWebTransportServer(&s.HTTP3.Server, cfg)
	if err != nil {
		return err
	}
	// Publish the server so the application can upgrade sessions from its
	// own request handler.
	webTransportMu.Lock()
	webTransportServer = wt
	webTransportMu.Unlock()
	// Regular h3 requests, including the WT CONNECT, go to the
	// application handler unchanged: the application decides which route
	// establishes a session.
	s.HTTP3.Handler = handler
	// webtransport.Server reuses the TLS config and QUICConfig of h3.
	return wt.ListenAndServe()
}
