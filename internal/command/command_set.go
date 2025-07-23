package command

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
	"strconv"
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

	var needTimeValue bool

	for _, arg := range arguments[2:] {
		if bulkText, ok := arg.(protocol.BulkString); ok {
			if needTimeValue {
				expiry, err := strconv.ParseInt(string(bulkText), 10, 64)
				if err != nil {
					return nil, NewSyntaxError()
				}
				cmd.expiry = expiry
				needTimeValue = false
				continue
			}

			switch string(bulkText) {
			case "GET":
				cmd.get = true
			case "NX":
				if cmd.existenceOption != existenceOptionNone {
					return nil, NewSyntaxError()
				}
				cmd.existenceOption = existenceOptionSetOnlyIfMissing
			case "XX":
				if cmd.existenceOption != existenceOptionNone {
					return nil, NewSyntaxError()
				}
				cmd.existenceOption = existenceOptionSetOnlyIfPresent
			case "EX":
				if cmd.expiryOption != expiryOptionNone {
					return nil, NewSyntaxError()
				}
				cmd.expiryOption = expiryOptionExpirySeconds
				needTimeValue = true
			case "PX":
				if cmd.expiryOption != expiryOptionNone {
					return nil, NewSyntaxError()
				}
				cmd.expiryOption = expiryOptionExpiryMilliseconds
				needTimeValue = true
			case "EXAT":
				if cmd.expiryOption != expiryOptionNone {
					return nil, NewSyntaxError()
				}
				cmd.expiryOption = expiryOptionExpiryUnixTimeInSeconds
				needTimeValue = true
			case "PXAT":
				if cmd.expiryOption != expiryOptionNone {
					return nil, NewSyntaxError()
				}
				cmd.expiryOption = expiryOptionExpiryUnixTimeInMilliseconds
				needTimeValue = true
			case "KEEPTTL":
				if cmd.expiryOption != expiryOptionNone {
					return nil, NewSyntaxError()
				}
				cmd.expiryOption = expiryOptionExpiryKeepTTL
			default:
				return nil, protocol.NewSimpleError("ERR wrong number of arguments for 'set' command")
			}
		} else {
			return nil, NewWrongDataTypeError(arg, protocol.BulkStringSymbol)
		}
	}

	if needTimeValue {
		return nil, NewSyntaxError()
	}

	return cmd, nil
}

type existenceOption string

const (
	existenceOptionSetOnlyIfMissing existenceOption = "NX"
	existenceOptionSetOnlyIfPresent existenceOption = "XX"
	existenceOptionNone             existenceOption = ""
)

type expiryOption string

const (
	expiryOptionNone                         expiryOption = ""
	expiryOptionExpirySeconds                expiryOption = "EX"
	expiryOptionExpiryMilliseconds           expiryOption = "PX"
	expiryOptionExpiryUnixTimeInSeconds      expiryOption = "EXAT"
	expiryOptionExpiryUnixTimeInMilliseconds expiryOption = "PXAT"
	expiryOptionExpiryKeepTTL                expiryOption = "KEEPTTL"
)

type SetCommand struct {
	key             string
	value           string
	get             bool
	existenceOption existenceOption
	expiryOption    expiryOption
	expiry          int64
}

func (cmd SetCommand) Execute(s store.Store) (protocol.Data, error) {
	oldValue, exists := s.Get(cmd.key)

	oldText, isText := parseStoreValueAsString(oldValue) // TODO validate setting against a list
	if !isText {
		return nil, nil
	}

	if exists && cmd.existenceOption == existenceOptionSetOnlyIfMissing {
		if cmd.get {
			return protocol.NewBulkString(oldText), nil
		}
		return nil, nil
	}
	if !exists && cmd.existenceOption == existenceOptionSetOnlyIfPresent {
		return nil, nil
	}

	if exists {
		swapped := s.CompareAndSwap(cmd.key, oldValue, cmd.value)
		if !swapped {
			return nil, nil
		}
	} else {
		_, loaded := s.LoadOrStore(cmd.key, cmd.value)
		if loaded {
			return nil, nil
		}
	}

	if cmd.get {
		if !exists {
			return nil, nil
		}
		return protocol.NewBulkString(oldText), nil
	}

	return protocol.NewSimpleString("OK"), nil
}

func parseStoreValueAsString(data any) (string, bool) {
	if data == nil {
		return "", true
	}
	if text, ok := data.(string); ok {
		return text, true
	}
	return "", false
}
