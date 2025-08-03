package command

import (
	"fmt"
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
)

type commandValidator interface {
	Validate(requestBytes []byte, arguments []protocol.Data) (Command, protocol.Data)
}

type Validator interface {
	Validate(requestBytes []byte, data protocol.Data) (Command, protocol.Data)
}

type validator struct {
	validators map[string]commandValidator
	clock      store.Clock
}

func NewValidator(clock store.Clock) Validator {
	return &validator{
		validators: map[string]commandValidator{
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
			"SET":    &SetValidator{clock: clock},
		},
		clock: clock,
	}
}

func (v *validator) Validate(requestBytes []byte, data protocol.Data) (Command, protocol.Data) {
	commandData, errorData := FromData(data)
	if errorData != nil {
		return nil, errorData
	}

	if selectedValidator, ok := v.validators[commandData.Name]; ok {
		return selectedValidator.Validate(requestBytes, commandData.Arguments)
	}

	return nil, protocol.NewSimpleError(fmt.Sprintf("ERR unknown command '%s'", commandData.Name))
}
