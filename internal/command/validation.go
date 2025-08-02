package command

import (
	"fmt"
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
)

var validators = map[string]commandValidator{
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
	"SET":    &SetValidator{clock: store.SystemClock{}},
}

type commandValidator interface {
	Validate(arguments []protocol.Data) (Command, protocol.Data)
}

type Validator interface {
	Validate(data protocol.Data) (Command, protocol.Data)
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
