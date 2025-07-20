package command

import (
	"redis-challenge/internal/protocol"
)

func validateSet(arguments []protocol.Data) (Command, protocol.Data) {
	if len(arguments) > 0 && arguments[0].Symbol() != protocol.BulkStringSymbol {
		return nil, NewWrongDataTypeError(arguments[0], protocol.BulkStringSymbol)
	}
	if len(arguments) > 1 && arguments[1].Symbol() != protocol.BulkStringSymbol {
		return nil, NewWrongDataTypeError(arguments[1], protocol.BulkStringSymbol)
	}

	if len(arguments) == 2 {
		return nil, nil
	}
	return nil, protocol.NewSimpleError("ERR wrong number of arguments for 'set' command")
}
