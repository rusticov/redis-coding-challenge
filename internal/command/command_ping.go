package command

import (
	"io"
	"redis-challenge/internal/protocol"
)

func validatePing(arguments []protocol.Data) (Command, protocol.Data) {
	if len(arguments) > 1 {
		return nil, protocol.NewSimpleError("ERR wrong number of arguments for 'ping' command")
	}
	return nil, nil
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
