package command

import (
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

func NewStoreExecutor(s store.Store, scanner Scanner) Executor {
	executionChannel := make(chan execution, callsToStoreQueueSize)

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			executionChannel <- execution{scan: scanner}
		}
	}()

	go func() { // TODO handle clean closing of this goroutine on server close
		for {
			e := <-executionChannel

			switch {
			case e.scan != nil:
				scanner.Scan()
			case e.cmd != nil:
				data, err := e.cmd.Execute(s)

				if err != nil {
					e.errors <- err
					continue
				}
				e.response <- data
			}
		}
	}()

	return storeExecutor{executionChannel: executionChannel}
}

type storeExecutor struct {
	executionChannel chan<- execution
}

func (executor storeExecutor) Execute(cmd Command, responses chan<- protocol.Data, errors chan<- error) {
	executor.executionChannel <- execution{cmd: cmd, errors: errors, response: responses}
}
