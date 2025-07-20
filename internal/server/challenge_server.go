package server

import (
	"bytes"
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
	defer func() {
		err := connection.Close()
		if err != nil {
			slog.Error("failed to close connection", "error", err)
		}
	}()

	var buffer bytes.Buffer

	readBuffer := make([]byte, 1024)
	for {
		bytesRead, err := connection.Read(readBuffer)
		if err != nil {
			slog.Error("failed to read request", "error", err)
			return
		}
		buffer.Write(readBuffer[:bytesRead])

		protocolData, requestByteCount := protocol.ReadFrame(buffer.Bytes())
		if requestByteCount == 0 {
			continue
		}
		data, commandError := command.FromData(protocolData)

		if commandError != nil {
			slog.Error("failed to parse request", "error", commandError, "request", buffer.String())
			err = protocol.WriteData(connection, commandError)
			if err != nil {
				slog.Error("failed to write parse request error", "error", err, "request", buffer.String())
			}
		} else {
			err = command.Registry{}.Execute(connection, data)
			if err != nil {
				slog.Error("failed to execute request", "error", err, "request", buffer.String())
			}
		}

		copy(readBuffer, buffer.Bytes()[requestByteCount:])
		buffer.Reset()
	}
}
