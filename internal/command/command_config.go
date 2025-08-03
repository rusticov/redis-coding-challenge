package command

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
)

type ConfigValidator struct{}

func (ConfigValidator) Validate(requestBytes []byte, _ []protocol.Data) (Command, protocol.Data) {
	return ConfigCommand{requestBytes: requestBytes}, nil
}

type ConfigCommand struct {
	requestBytes []byte
}

func (cmd ConfigCommand) Request() ([]byte, Type) {
	return cmd.requestBytes, TypeRead
}

func (cmd ConfigCommand) Execute(_ store.Store) (protocol.Data, error) {
	return nil, nil
}
