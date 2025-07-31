package command

import (
	"redis-challenge/internal/protocol"
)

type IncrValidator struct{}

func (IncrValidator) Validate(arguments []protocol.Data) (Command, protocol.Data) {
	var key string

	if len(arguments) > 0 {
		if arg, ok := arguments[0].(protocol.BulkString); ok {
			key = string(arg)
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

	return ChangeIntegerCommand{key: key, change: 1}, nil
}
