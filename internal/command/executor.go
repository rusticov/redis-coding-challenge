package command

import (
	"context"
	"io"
	"log/slog"
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
	"time"
)

type Executor interface {
	Execute(request []byte, cmd Command, responses chan<- protocol.Data, errors chan<- error)
}

type Scanner interface {
	Scan()
}

type execution struct {
	cmd      Command
	request  []byte
	scan     Scanner
	errors   chan<- error
	response chan<- protocol.Data
	writer   io.Writer
}

const callsToStoreQueueSize = 1000

func NewStoreExecutor(ctx context.Context, s store.Store, scanner Scanner, writer io.Writer) Executor {
	executionChannel := make(chan execution, callsToStoreQueueSize)

	go triggerRepeatedExpiryScan(ctx, executionChannel, scanner)

	go executeCommandsAgainstStore(ctx, executionChannel, s, writer)

	return storeExecutor{executionChannel: executionChannel}
}

func triggerRepeatedExpiryScan(ctx context.Context, executionChannel chan<- execution, scanner Scanner) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-ctx.Done():
			return
		default:
			executionChannel <- execution{scan: scanner}
		}
	}
}

func executeCommandsAgainstStore(ctx context.Context, executionChannel <-chan execution, s store.Store, writer io.Writer) {
	for {
		select {
		case <-ctx.Done():
			return
		case e := <-executionChannel:
			switch {
			case e.scan != nil:
				e.scan.Scan()
			case e.cmd != nil:
				if e.cmd.IsUpdate() {
					_, err := writer.Write(e.request)
					if err != nil {
						slog.Error("failed to write request", "error", err, "request", string(e.request))
						return
					}
				}

				data, err := e.cmd.Execute(s)
				if err != nil {
					e.errors <- err
					continue
				}

				e.response <- data
			}
		}
	}
}

type storeExecutor struct {
	executionChannel chan<- execution
}

func (executor storeExecutor) Execute(request []byte, cmd Command, responses chan<- protocol.Data, errors chan<- error) {
	executor.executionChannel <- execution{cmd: cmd, request: request, errors: errors, response: responses}
}
