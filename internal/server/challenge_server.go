package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"redis-challenge/internal/command"
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
)

type State int

const (
	StatePortClosed State = iota + 1
)

type MonitorChannel chan State

type ChallengeServer struct {
	socket         net.Listener
	cancelFunction context.CancelFunc
	monitor        MonitorChannel
}

func (c *ChallengeServer) Address() string {
	return c.socket.Addr().String()
}

func (c *ChallengeServer) Close() error {
	c.cancelFunction()
	err := c.socket.Close()

	if c.monitor != nil {
		c.monitor <- StatePortClosed
	}

	return err
}

func (c *ChallengeServer) AddMonitor(monitor MonitorChannel) {
	c.monitor = monitor
}

func NewChallengeServer(port int, s store.Store, scanner command.Scanner) (*ChallengeServer, error) {
	executor := command.NewStoreExecutor(s, scanner)

	socket, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
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
					if errors.Is(err, net.ErrClosed) {
						return
					}
					slog.Error("failed to accept on socket", "error", err)
					continue
				}

				go func() {
					connectionHandler(connection, executor)
				}()
			}
		}
	}()

	return &ChallengeServer{
		socket:         socket,
		cancelFunction: cancelFunction,
	}, nil
}

func connectionHandler(connection net.Conn, executor command.Executor) {
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
			if err == io.EOF {
				return
			}

			slog.Error("failed to read request", "error", err)
			return
		}
		buffer.Write(readBuffer[:bytesRead])

		protocolData, requestByteCount := protocol.ReadFrame(buffer.Bytes())
		if requestByteCount == 0 {
			continue
		}
		response := executeCommand(protocolData, buffer, executor)

		err = protocol.WriteData(connection, response)
		if err != nil {
			slog.Error("failed to write parse request error", "error", err, "request", buffer.String())
		}

		copy(readBuffer, buffer.Bytes()[requestByteCount:])
		buffer.Reset()
	}
}

func executeCommand(protocolData protocol.Data, buffer bytes.Buffer, executor command.Executor) protocol.Data {
	parsedCommand, commandError := command.Validate(protocolData)

	switch {
	case commandError != nil:
		slog.Error("failed to parse request", "error", commandError, "request", buffer.String())
		return commandError
	case parsedCommand == nil:
		slog.Error("expect a command if there is no error data on parsing", "error", commandError, "request", buffer.String())
		return protocol.NewSimpleError("ERR protocol error")
	default:
		responseReceiver := make(chan protocol.Data)
		errorReceiver := make(chan error)

		executor.Execute(parsedCommand, responseReceiver, errorReceiver)

		select {
		case err := <-errorReceiver:
			slog.Error("failed to execute request", "error", err, "request", buffer.String())
			return protocol.NewSimpleError("ERR protocol error")

		case response := <-responseReceiver:
			return response
		}
	}
}
