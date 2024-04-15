package server

// TODO
type ServerOption interface {
	apply(server *Server)
}
