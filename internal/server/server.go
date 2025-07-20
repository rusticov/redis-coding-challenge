package server

// Server represents a running instance of a Redis-compliant server.
type Server interface {
	Address() string
	Close() error
}
