package command

import (
	"errors"
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
)

type RPushValidator struct{}

func (RPushValidator) Validate(arguments []protocol.Data) (Command, protocol.Data) {
	values := make([]string, len(arguments))
	for i, arg := range arguments {
		if _, ok := arguments[i].(protocol.BulkString); ok {
			values[i] = string(arg.(protocol.BulkString))
			continue
		}

		return nil, NewWrongDataTypeError(arguments[i], protocol.BulkStringSymbol)
	}

	if len(arguments) < 2 {
		return nil, protocol.NewSimpleError("ERR wrong number of arguments for 'rpush' command")
	}

	return RPushCommand{
		key:    values[0],
		values: values[1:],
	}, nil
}

type RPushCommand struct {
	key    string
	values []string
}

func (cmd RPushCommand) Execute(s store.Store) (protocol.Data, error) {
	count, err := s.RightPush(cmd.key, cmd.values)

	if errors.Is(err, store.ErrorWrongOperationType) {
		return NewWrongOperationTypeError(), nil
	}
	if err != nil {
		return nil, err
	}

	return protocol.NewSimpleInteger(count), nil
}
