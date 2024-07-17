package cherry

// TODO
type ServerOption interface {
	apply(server *Server)
}
