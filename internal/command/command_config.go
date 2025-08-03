package command

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
)

type ConfigValidator struct{}

func (ConfigValidator) Validate(_ []protocol.Data) (Command, protocol.Data) {
	return ConfigCommand{}, nil
}

type ConfigCommand struct {
}

func (cmd ConfigCommand) IsUpdate() bool {
	return false
}

func (cmd ConfigCommand) Execute(_ store.Store) (protocol.Data, error) {
	return nil, nil
}
