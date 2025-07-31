package command

import (
	"fmt"
	"redis-challenge/internal/protocol"
)

var validators = map[string]validator{
	"PING":   PingValidator{},
	"ECHO":   EchoValidator{},
	"CONFIG": ConfigValidator{},
	"DECR":   DecrValidator{},
	"DEL":    DelValidator{},
	"EXISTS": ExistsValidator{},
	"INCR":   IncrValidator{},
	"GET":    GetValidator{},
	"LPUSH":  LPushValidator{},
	"LRANGE": LRangeValidator{},
	"RPUSH":  RPushValidator{},
	"SET":    SetValidator{},
}

type validator interface {
	Validate(arguments []protocol.Data) (Command, protocol.Data)
}

func Validate(data protocol.Data) (Command, protocol.Data) {
	commandData, errorData := FromData(data)
	if errorData != nil {
		return nil, errorData
	}

	if v, ok := validators[commandData.Name]; ok {
		return v.Validate(commandData.Arguments)
	}

	return nil, protocol.NewSimpleError(fmt.Sprintf("ERR unknown command '%s'", commandData.Name))
}
