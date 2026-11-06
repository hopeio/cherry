package mix

import "testing"

// prepareHTTP3 owns the HTTP/3 defaults and the "no TLS material" outcome, so both are
// pinned here rather than left to whatever each consumer re-implements.

func TestPrepareHTTP3CertFallsBackToServer(t *testing.T) {
	s := &Server{}
	s.CertFile, s.KeyFile = "server.crt", "server.key"
	s.HTTP3.Enabled = true

	s.prepareHTTP3()

	if !s.HTTP3.Enabled {
		t.Fatal("h3 must stay enabled when the main server's certificate can be reused")
	}
	if s.HTTP3.CertFile != "server.crt" || s.HTTP3.KeyFile != "server.key" {
		t.Fatalf("cert pair not inherited: %q/%q", s.HTTP3.CertFile, s.HTTP3.KeyFile)
	}
}

func TestPrepareHTTP3ExplicitCertWins(t *testing.T) {
	s := &Server{}
	s.CertFile, s.KeyFile = "server.crt", "server.key"
	s.HTTP3.CertFile, s.HTTP3.KeyFile = "h3.crt", "h3.key"
	s.HTTP3.Enabled = true

	s.prepareHTTP3()

	if s.HTTP3.CertFile != "h3.crt" || s.HTTP3.KeyFile != "h3.key" {
		t.Fatalf("explicit h3 pair overwritten: %q/%q", s.HTTP3.CertFile, s.HTTP3.KeyFile)
	}
}

func TestPrepareHTTP3MissingPairDisablesInsteadOfFailing(t *testing.T) {
	s := &Server{}
	s.HTTP3.Enabled = true

	s.prepareHTTP3()

	if s.HTTP3.Enabled {
		t.Fatal("a QUIC listener without TLS must disable itself")
	}
	if s.Addr != "" || s.CertFile != "" {
		t.Fatal("preparation must not touch the main server")
	}
}

func TestPrepareHTTP3AddrDefault(t *testing.T) {
	s := &Server{}
	s.CertFile, s.KeyFile = "server.crt", "server.key"
	s.HTTP3.Enabled = true
	s.prepareHTTP3()
	if s.HTTP3.Addr != ":443" {
		t.Fatalf("h3 Addr=%q, want :443 (quic-go binds :https, the log must name the open port)", s.HTTP3.Addr)
	}

	s2 := &Server{}
	s2.CertFile, s2.KeyFile = "server.crt", "server.key"
	s2.HTTP3.Enabled = true
	s2.HTTP3.Addr = ":9443"
	s2.prepareHTTP3()
	if s2.HTTP3.Addr != ":9443" {
		t.Fatalf("explicit h3 Addr overwritten: %q", s2.HTTP3.Addr)
	}
}

func TestPrepareHTTP3DisabledStaysDisabled(t *testing.T) {
	s := &Server{}
	s.CertFile, s.KeyFile = "server.crt", "server.key"

	s.prepareHTTP3()

	if s.HTTP3.Enabled {
		t.Fatal("preparation only ever disables, never enables")
	}
	if s.HTTP3.CertFile != "" || s.HTTP3.Addr != "" {
		t.Fatalf("work done on a disabled listener: %+v", s.HTTP3)
	}
}
