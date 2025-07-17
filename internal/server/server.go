package server

// Server represents a running instance of a Redis-compliant server.
type Server interface {
	Address() string
	Close() error
}

// RealRedisServer exists to running to a real running Redis server for comparison with the one written for this challenge.
type RealRedisServer struct{}

func (s RealRedisServer) Address() string {
	return ":6379"
}

func (s RealRedisServer) Close() error {
	return nil
}

// NewRealRedisServer initializes and returns a new instance of RealRedisServer, representing a real Redis server.
func NewRealRedisServer() Server {
	return RealRedisServer{}
}
