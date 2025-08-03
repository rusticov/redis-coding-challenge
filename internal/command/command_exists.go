package command

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
)

type ExistsValidator struct{}

func (ExistsValidator) Validate(requestBytes []byte, arguments []protocol.Data) (Command, protocol.Data) {
	if len(arguments) == 0 {
		return nil, protocol.NewSimpleError("ERR wrong number of arguments for 'exists' command")
	}

	keys := make([]string, len(arguments))
	for i, arg := range arguments {
		if _, ok := arguments[i].(protocol.BulkString); ok {
			keys[i] = string(arg.(protocol.BulkString))
			continue
		}

		return nil, NewWrongDataTypeError(arguments[i], protocol.BulkStringSymbol)
	}

	return ExistsCommand{
		requestBytes: requestBytes,
		keys:         keys,
	}, nil
}

type ExistsCommand struct {
	requestBytes []byte
	keys         []string
}

func (cmd ExistsCommand) Request() ([]byte, Type) {
	return cmd.requestBytes, TypeRead
}

func (cmd ExistsCommand) Execute(s store.Store) (protocol.Data, error) {
	count := 0
	for _, key := range cmd.keys {
		if s.Exists(key) {
			count++
		}
	}
	return protocol.NewSimpleInteger(int64(count)), nil
}
