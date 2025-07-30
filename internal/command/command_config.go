package command

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
)

func validateConfig(_ []protocol.Data) (Command, protocol.Data) {
	return ConfigCommand{}, nil
}

type ConfigCommand struct {
}

func (cmd ConfigCommand) Execute(s store.Store) (protocol.Data, error) {
	return nil, nil
}
