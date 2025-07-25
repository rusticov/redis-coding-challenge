package command

import (
	"errors"
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
	"redis-challenge/internal/store/list"
)

func validateLPush(arguments []protocol.Data) (Command, protocol.Data) {
	values := make([]string, len(arguments))
	for i, arg := range arguments {
		if _, ok := arguments[i].(protocol.BulkString); ok {
			values[i] = string(arg.(protocol.BulkString))
			continue
		}

		return nil, NewWrongDataTypeError(arguments[i], protocol.BulkStringSymbol)
	}

	if len(arguments) < 2 {
		return nil, protocol.NewSimpleError("ERR wrong number of arguments for 'lpush' command")
	}

	return LPushCommand{
		key:    values[0],
		values: values[1:],
	}, nil
}

type LPushCommand struct {
	key    string
	values []string
}

func (cmd LPushCommand) Execute(s store.Store) (protocol.Data, error) {
	count, err := s.LeftPush(cmd.key, cmd.values)

	if errors.Is(err, list.ErrorOldValueIsNotList) {
		return NewWrongOperationTypeError(), nil
	}
	if err != nil {
		return nil, err
	}

	return protocol.NewSimpleInteger(count), nil
}
