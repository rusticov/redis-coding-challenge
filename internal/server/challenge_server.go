package server

import (
	"log/slog"
	"net"
)

type ChallengeServer struct {
	socket net.Listener
}

func (c *ChallengeServer) Address() string {
	return c.socket.Addr().String()
}

func (c *ChallengeServer) Close() error {
	return c.socket.Close()
}

func NewChallengeServer() (Server, error) {
	socket, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			connection, err := socket.Accept()
			if err != nil {
				slog.Error("failed to accept connection", "error", err)
				continue
			}

			go func() {
				commandHandler(connection)
			}()
		}
	}()

	return &ChallengeServer{socket: socket}, nil
}

func commandHandler(connection net.Conn) {
	defer connection.Close()

	connection.Write([]byte("+PONG\r\n"))
}
