package command

import (
	"fmt"
	"redis-challenge/internal/protocol"
)

var validators = map[string]validator{
	"PING":   validatePing,
	"ECHO":   validateEcho,
	"DECR":   validateDecr,
	"DEL":    validateDel,
	"EXISTS": validateExists,
	"INCR":   validateIncr,
	"GET":    validateGet,
	"SET":    validateSet,
}

type validator func(arguments []protocol.Data) (Command, protocol.Data)

func Validate(data protocol.Data) (Command, protocol.Data) {
	commandData, errorData := FromData(data)
	if errorData != nil {
		return nil, errorData
	}

	if v, ok := validators[commandData.Name]; ok {
		return v(commandData.Arguments)
	}

	return nil, protocol.NewSimpleError(fmt.Sprintf("ERR unknown command '%s'", commandData.Name))
}
