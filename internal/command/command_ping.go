package command

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
)

type PingValidator struct{}

func (PingValidator) Validate(requestBytes []byte, arguments []protocol.Data) (Command, protocol.Data) {
	switch len(arguments) {
	case 0:
		return PingCommand{
			requestBytes: requestBytes,
			response:     protocol.NewSimpleString("PONG"),
		}, nil
	case 1:
		if _, ok := arguments[0].(protocol.BulkString); !ok {
			return nil, NewWrongDataTypeError(arguments[0], protocol.BulkStringSymbol)
		}
		return PingCommand{
			requestBytes: requestBytes,
			response:     arguments[0],
		}, nil
	default:
		return nil, protocol.NewSimpleError("ERR wrong number of arguments for 'ping' command")
	}
}

type PingCommand struct {
	requestBytes []byte
	response     protocol.Data
}

func (cmd PingCommand) Request() ([]byte, Type) {
	return cmd.requestBytes, TypeRead
}

func (cmd PingCommand) Execute(_ store.Store) (protocol.Data, error) {
	return cmd.response, nil
}
