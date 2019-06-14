package server

import "net/http"

// Option represents the Server options
type Option func(s *Server)

// WithHandler sets the handler
func WithHandler(handler http.Handler) Option {
	return func(s *Server) {
		s.handler = handler
	}
}
