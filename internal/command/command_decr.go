package command

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
	"strconv"
)

func validateDecr(arguments []protocol.Data) (Command, protocol.Data) {
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
		return nil, protocol.NewSimpleError("ERR wrong number of arguments for 'decr' command")
	}

	return DecrCommand{key: key}, nil
}

type DecrCommand struct {
	key string
}

func (cmd DecrCommand) Execute(s store.Store) (protocol.Data, error) {
	textValue, exists := s.Get(cmd.key)

	if !exists {
		s.LoadOrStore(cmd.key, "-1") // TODO test failure to set here
		return protocol.NewSimpleInteger(-1), nil
	}

	var value int64
	if textValue != "" {
		var err error
		value, err = strconv.ParseInt(textValue, 10, 64)
		if err != nil {
			return protocol.NewSimpleError("ERR value is not an integer or out of range"), nil
		}
	}
	value--

	s.CompareAndSwap(cmd.key, textValue, strconv.FormatInt(value, 10)) // TODO test failure to set here

	return protocol.NewSimpleInteger(value), nil
}
