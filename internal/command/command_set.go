package command

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
)

func validateSet(arguments []protocol.Data) (Command, protocol.Data) {
	if len(arguments) > 0 && arguments[0].Symbol() != protocol.BulkStringSymbol {
		return nil, NewWrongDataTypeError(arguments[0], protocol.BulkStringSymbol)
	}
	if len(arguments) > 1 && arguments[1].Symbol() != protocol.BulkStringSymbol {
		return nil, NewWrongDataTypeError(arguments[1], protocol.BulkStringSymbol)
	}

	if len(arguments) < 2 {
		return nil, protocol.NewSimpleError("ERR wrong number of arguments for 'set' command")
	}

	cmd := SetCommand{
		key:   string(arguments[0].(protocol.BulkString)),
		value: string(arguments[1].(protocol.BulkString)),
	}

	for _, arg := range arguments[2:] {
		if bulkText, ok := arg.(protocol.BulkString); ok {
			switch string(bulkText) {
			case "GET":
				cmd.get = true
			case "NX":
				if cmd.existenceOption != valueExistenceOptionNone {
					return nil, NewSyntaxError()
				}
				cmd.existenceOption = valueExistenceOptionPresent
			case "XX":
				if cmd.existenceOption != valueExistenceOptionNone {
					return nil, NewSyntaxError()
				}
				cmd.existenceOption = valueExistenceOptionAbsent
			default:
				return nil, protocol.NewSimpleError("ERR wrong number of arguments for 'set' command")
			}
		} else {
			return nil, NewWrongDataTypeError(arg, protocol.BulkStringSymbol)
		}
	}

	return cmd, nil
}

type valueExistenceOption string

const (
	valueExistenceOptionPresent valueExistenceOption = "NX"
	valueExistenceOptionAbsent  valueExistenceOption = "XX"
	valueExistenceOptionNone    valueExistenceOption = ""
)

type SetCommand struct {
	key             string
	value           string
	get             bool
	existenceOption valueExistenceOption
}

func (cmd SetCommand) Execute(s *store.Store) (protocol.Data, error) {
	oldValue, exists := s.Get(cmd.key)
	s.Add(cmd.key, cmd.value)

	if cmd.get {
		if !exists {
			return nil, nil
		}
		return protocol.NewBulkString(oldValue), nil
	}

	return protocol.NewSimpleString("OK"), nil
}
