package server

type Server struct{}

func (s *Server) Address() string {
	return ":6379"
}

func NewServer() (*Server, error) {
	return nil, nil
}
