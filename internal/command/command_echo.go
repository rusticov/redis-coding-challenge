package command

import (
	"io"
	"redis-challenge/internal/protocol"
)

type EchoCommand struct {
}

func (EchoCommand) Execute(writer io.Writer, data Data) error {
	if len(data.Arguments) == 0 {
		return protocol.WriteData(writer, protocol.NewSimpleError("ERR wrong number of arguments for 'echo' command"))
	}
	return protocol.WriteData(writer, data.Arguments[0])
}
