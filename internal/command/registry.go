package command

import (
	"fmt"
	"io"
	"redis-challenge/internal/protocol"
)

var commands = map[string]Command{
	"PING": PingCommand{},
	"ECHO": EchoCommand{},
}

type Command interface {
	Execute(writer io.Writer, data Data) error
}

type Registry struct {
}

func (r Registry) Execute(writer io.Writer, data Data) error {
	if command, ok := commands[data.Name]; ok {
		return command.Execute(writer, data)
	}

	errorMessage := fmt.Sprintf("ERR unknown command '%s'", data.Name)

	return protocol.WriteData(writer, protocol.NewSimpleError(errorMessage))
}
