package command

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
)

func validateSet(arguments []protocol.Data) (Command, protocol.Data) {
	if len(arguments) > 0 && arguments[0].Symbol() != protocol.BulkStringSymbol {
		return nil, NewWrongDataTypeError(arguments[0], protocol.BulkStringSymbol)
	}
	if len(arguments) > 1 && arguments[1].Symbol() != protocol.BulkStringSymbol {
		return nil, NewWrongDataTypeError(arguments[1], protocol.BulkStringSymbol)
	}

	if len(arguments) == 2 {
		return SetCommand{
			key:   string(arguments[0].(protocol.BulkString)),
			value: string(arguments[1].(protocol.BulkString)),
		}, nil
	}
	return nil, protocol.NewSimpleError("ERR wrong number of arguments for 'set' command")
}

type SetCommand struct {
	key   string
	value string
}

func (cmd SetCommand) Execute(s *store.Store) (protocol.Data, error) {
	s.Add(cmd.key, cmd.value)
	return protocol.NewSimpleString("OK"), nil
}
