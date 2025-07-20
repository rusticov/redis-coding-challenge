package command

import (
	"redis-challenge/internal/protocol"
)

func validateGet(arguments []protocol.Data) (Command, protocol.Data) {
	if len(arguments) != 1 {
		return nil, protocol.NewSimpleError("ERR wrong number of arguments for 'get' command")
	}
	if _, ok := arguments[0].(protocol.BulkString); !ok {
		return nil, NewWrongDataTypeError(arguments[0], protocol.BulkStringSymbol)
	}
	return nil, nil
}
