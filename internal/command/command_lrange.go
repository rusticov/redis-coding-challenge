package command

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
	"strconv"
)

func validateLRange(arguments []protocol.Data) (Command, protocol.Data) {
	var cmd LRangeCommand

	if len(arguments) > 0 {
		if arg, ok := arguments[0].(protocol.BulkString); ok {
			cmd.key = string(arg)
		} else {
			return nil, NewWrongDataTypeError(arguments[0], protocol.BulkStringSymbol)
		}
	}

	if len(arguments) > 1 {
		value, errorData := parseIntegerFromData(arguments[1])
		if errorData != nil {
			return nil, errorData
		}
		cmd.left = value
	}

	if len(arguments) > 2 {
		value, errorData := parseIntegerFromData(arguments[2])
		if errorData != nil {
			return nil, errorData
		}
		cmd.right = value
	}

	if len(arguments) != 3 {
		return nil, protocol.NewSimpleError("ERR wrong number of arguments for 'lrange' command")
	}

	return cmd, nil
}

func parseIntegerFromData(data protocol.Data) (int, protocol.Data) {
	if arg, ok := data.(protocol.BulkString); ok {
		intValue, err := strconv.Atoi(string(arg))
		if err != nil {
			return 0, protocol.NewSimpleError("ERR value is not an integer or out of range")
		}
		return intValue, nil
	} else {
		return 0, NewWrongDataTypeError(data, protocol.BulkStringSymbol)
	}
}

type LRangeCommand struct {
	key   string
	left  int
	right int
}

func (cmd LRangeCommand) Execute(s store.Store) (protocol.Data, error) {
	return protocol.NewSimpleInteger(0), nil
}
