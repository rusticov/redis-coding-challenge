package server

import (
	"context"
	"log/slog"
	"net"
	"redis-challenge/internal/command"
	"redis-challenge/internal/protocol"
)

type ChallengeServer struct {
	socket         net.Listener
	cancelFunction context.CancelFunc
}

func (c *ChallengeServer) Address() string {
	return c.socket.Addr().String()
}

func (c *ChallengeServer) Close() error {
	c.cancelFunction()
	return c.socket.Close()
}

func NewChallengeServer() (Server, error) {
	socket, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, err
	}

	ctx, cancelFunction := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				connection, err := socket.Accept()
				if err != nil {
					slog.Error("failed to accept connection", "error", err)
					continue
				}

				go func() {
					connectionHandler(connection)
				}()
			}
		}
	}()

	return &ChallengeServer{
		socket:         socket,
		cancelFunction: cancelFunction,
	}, nil
}

func connectionHandler(connection net.Conn) {
	defer connection.Close()

	buffer := make([]byte, 1024)
	n, err := connection.Read(buffer)
	if err != nil {
		// TODO return error message
		return
	}

	protocolData, _ := protocol.ReadFrame(buffer[:n])
	data, _ := command.FromData(protocolData) // TODO respond with error data

	command.Registry{}.Execute(connection, data) // TODO handle error
}
