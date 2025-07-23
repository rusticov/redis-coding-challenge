package command

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
)

func validateIncr(arguments []protocol.Data) (Command, protocol.Data) {
	if len(arguments) > 0 {
		if _, ok := arguments[0].(protocol.BulkString); ok {
			//
		} else {
			return nil, NewWrongDataTypeError(arguments[0], protocol.BulkStringSymbol)
		}
	}

	if len(arguments) > 1 {
		_, ok := arguments[1].(protocol.BulkString)
		if !ok {
			return nil, NewWrongDataTypeError(arguments[1], protocol.BulkStringSymbol)
		}
	}

	if len(arguments) != 1 {
		return nil, protocol.NewSimpleError("ERR wrong number of arguments for 'incr' command")
	}

	return IncrCommand{}, nil
}

type IncrCommand struct {
}

func (cmd IncrCommand) Execute(s store.Store) (protocol.Data, error) {
	return nil, nil
}
