package command

import (
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

func (cmd GetCommand) Execute(s *store.Store) (protocol.Data, error) {
	value, exists := s.Get(cmd.key)
	if !exists {
		return nil, nil
	}
	return protocol.NewBulkString(value), nil
}
