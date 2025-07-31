package command

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
)

type EchoValidator struct{}

func (EchoValidator) Validate(arguments []protocol.Data) (Command, protocol.Data) {
	if len(arguments) != 1 {
		return nil, protocol.NewSimpleError("ERR wrong number of arguments for 'echo' command")
	}
	if _, ok := arguments[0].(protocol.BulkString); !ok {
		return nil, NewWrongDataTypeError(arguments[0], protocol.BulkStringSymbol)
	}
	return EchoCommand{response: arguments[0]}, nil
}

type EchoCommand struct {
	response protocol.Data
}

func (cmd EchoCommand) Execute(_ store.Store) (protocol.Data, error) {
	return cmd.response, nil
}
