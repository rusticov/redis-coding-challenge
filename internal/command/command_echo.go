package command

import (
	"io"
	"redis-challenge/internal/protocol"
)

func validateEcho(arguments []protocol.Data) (Command, protocol.Data) {
	if len(arguments) != 1 {
		return nil, protocol.NewSimpleError("ERR wrong number of arguments for 'echo' command")
	}
	if _, ok := arguments[0].(protocol.BulkString); !ok {
		return nil, NewWrongDataTypeError(arguments[0], protocol.BulkStringSymbol)
	}
	return nil, nil
}

type EchoCommand struct {
}

func (EchoCommand) Execute(writer io.Writer, data Data) error {
	if len(data.Arguments) != 1 {
		return protocol.WriteData(writer, protocol.NewSimpleError("ERR wrong number of arguments for 'echo' command"))
	}
	return protocol.WriteData(writer, data.Arguments[0])
}
