package server

// Server represents a running instance of a Redis-compliant server.
type Server interface {
	Address() string
}

// RealRedisServer exists to running to a real running Redis server for comparison with the one written for this challenge.
type RealRedisServer struct{}

func (s *RealRedisServer) Address() string {
	return ":6379"
}

// NewRealRedisServer initializes and returns a new instance of RealRedisServer, representing a real Redis server.
func NewRealRedisServer() (Server, error) {
	return &RealRedisServer{}, nil
}
