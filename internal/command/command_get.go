package command

import (
	"errors"
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
)

func validateGet(arguments []protocol.Data) (Command, protocol.Data) {
	if len(arguments) != 1 {
		return nil, protocol.NewSimpleError("ERR wrong number of arguments for 'get' command")
	}
	if _, ok := arguments[0].(protocol.BulkString); !ok {
		return nil, NewWrongDataTypeError(arguments[0], protocol.BulkStringSymbol)
	}
	return GetCommand{
		key: string(arguments[0].(protocol.BulkString)),
	}, nil
}

type GetCommand struct {
	key string
}

func (cmd GetCommand) Execute(s store.Store) (protocol.Data, error) {
	value, err := s.ReadString(cmd.key)

	if errors.Is(err, store.ErrorKeyNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return protocol.NewBulkString(value), err
}
