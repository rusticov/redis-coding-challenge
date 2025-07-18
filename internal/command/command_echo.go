package command

import (
	"io"
	"redis-challenge/internal/protocol"
)

type EchoCommand struct {
}

func (EchoCommand) Execute(writer io.Writer, data Data) error {
	return protocol.WriteData(writer, data.Arguments[0])
}
