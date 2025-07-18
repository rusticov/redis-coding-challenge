package command

import (
	"io"
	"redis-challenge/internal/protocol"
)

type PingCommand struct {
}

func (PingCommand) Execute(writer io.Writer, data Data) error {
	if len(data.Arguments) == 0 {
		_, err := writer.Write([]byte("+PONG\r\n"))
		return err
	}

	return protocol.WriteData(writer, data.Arguments[0])
}
