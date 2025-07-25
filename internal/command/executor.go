package command

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
)

type Executor interface {
	Execute(cmd Command, responses chan<- protocol.Data, errors chan<- error)
}

type execution struct {
	cmd      Command
	errors   chan<- error
	response chan<- protocol.Data
}

func NewStoreExecutor(s store.Store) Executor {
	executionChannel := make(chan execution, 1000)

	go func() {
		for {
			e := <-executionChannel

			data, err := e.cmd.Execute(s)

			if err != nil {
				e.errors <- err
				continue
			}
			e.response <- data
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
