package command

import (
	"io"
	"redis-challenge/internal/protocol"
)

func validatePing(arguments []protocol.Data) (Command, protocol.Data) {
	switch len(arguments) {
	case 0:
		return nil, nil
	case 1:
		if _, ok := arguments[0].(protocol.BulkString); !ok {
			return nil, NewWrongDataTypeError(arguments[0], protocol.BulkStringSymbol)
		}
		return nil, nil
	default:
		return nil, protocol.NewSimpleError("ERR wrong number of arguments for 'ping' command")
	}
}

type PingCommand struct {
}

func (PingCommand) Execute(writer io.Writer, data Data) error {
	if len(data.Arguments) == 0 {
		_, err := writer.Write([]byte("+PONG\r\n"))
		return err
	}

	return protocol.WriteData(writer, data.Arguments[0])
}
