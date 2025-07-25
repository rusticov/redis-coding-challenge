package command

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
)

type Executor interface {
	Execute(cmd Command, responses chan<- protocol.Data, errors chan<- error)
}

func NewStoreExecutor(store store.Store) Executor {
	return StoreExecutor{store: store}
}

type StoreExecutor struct {
	store store.Store
}

func (executor StoreExecutor) Execute(cmd Command, responses chan<- protocol.Data, errors chan<- error) {
	go func() {
		data, err := cmd.Execute(executor.store)
		if err != nil {
			errors <- err
			return
		}

		responses <- data
	}()
}
