package command

import (
	"context"
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
	"time"
)

type Executor interface {
	Execute(cmd Command, responses chan<- protocol.Data, errors chan<- error)
}

type Scanner interface {
	Scan()
}

type execution struct {
	cmd      Command
	scan     Scanner
	errors   chan<- error
	response chan<- protocol.Data
}

const callsToStoreQueueSize = 1000

func NewStoreExecutor(ctx context.Context, s store.Store, scanner Scanner) Executor {
	executionChannel := make(chan execution, callsToStoreQueueSize)

	go triggerRepeatedExpiryScan(ctx, executionChannel, scanner)

	go executeCommandsAgainstStore(ctx, executionChannel, s)

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

func executeCommandsAgainstStore(ctx context.Context, executionChannel <-chan execution, s store.Store) {
	for {
		select {
		case <-ctx.Done():
			return
		case e := <-executionChannel:
			switch {
			case e.scan != nil:
				e.scan.Scan()
			case e.cmd != nil:
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

func (executor storeExecutor) Execute(cmd Command, responses chan<- protocol.Data, errors chan<- error) {
	executor.executionChannel <- execution{cmd: cmd, errors: errors, response: responses}
}
