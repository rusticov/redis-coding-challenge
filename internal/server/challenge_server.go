package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"

	"redis-challenge/internal/command"
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

type ChallengeServerBuilder struct {
	port           int
	builder        store.Builder
	writer         io.Writer
	reader         io.Reader
	monitorChannel MonitorChannel
	err            error
}

func NewChallengeServer(port int, builder store.Builder) *ChallengeServerBuilder {
	return &ChallengeServerBuilder{
		port:    port,
		builder: builder,
		writer:  io.Discard,
	}
}

func (b *ChallengeServerBuilder) WithArchiveWriter(writer io.Writer) *ChallengeServerBuilder {
	b.writer = writer
	return b
}

func (b *ChallengeServerBuilder) WithMonitorChannel(monitorChannel MonitorChannel) *ChallengeServerBuilder {
	b.monitorChannel = monitorChannel
	return b
}

func (b *ChallengeServerBuilder) Start() (*ChallengeServer, error) {
	if b.err != nil {
		return nil, b.err
	}

	socket, err := net.Listen("tcp", fmt.Sprintf(":%d", b.port))
	if err != nil {
		return nil, err
	}

	s, scanner := b.builder.WithCommandLogWriter(b.writer).Build()

	if b.reader != nil {
		r := restorer{store: s}

		err = r.RestoreFromLog(b.reader)
		if err != nil {
			return nil, fmt.Errorf("failed restore from log: %w", err)
		}
	}

	ctx, cancelFunction := context.WithCancel(context.Background())

	handler := connectionHandler{
		executor: command.NewStoreExecutor(ctx, s, scanner, b.writer),
	}

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
					handler.HandleConnection(connection)
				}()
			}
		}
	}()

	return &ChallengeServer{
		socket:         socket,
		cancelFunction: cancelFunction,
	}, nil
}

func (b *ChallengeServerBuilder) RestoreFromArchive(reader io.Reader) *ChallengeServerBuilder {
	b.reader = reader
	return b
}
