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
				if cmd.expiryOption != store.ExpiryOptionNone {
					return nil, NewSyntaxError()
				}
				cmd.expiryOption = store.ExpiryOptionExpirySeconds
				needTimeValue = true
			case "PX":
				if cmd.expiryOption != store.ExpiryOptionNone {
					return nil, NewSyntaxError()
				}
				cmd.expiryOption = store.ExpiryOptionExpiryMilliseconds
				needTimeValue = true
			case "EXAT":
				if cmd.expiryOption != store.ExpiryOptionNone {
					return nil, NewSyntaxError()
				}
				cmd.expiryOption = store.ExpiryOptionExpiryUnixTimeInSeconds
				needTimeValue = true
			case "PXAT":
				if cmd.expiryOption != store.ExpiryOptionNone {
					return nil, NewSyntaxError()
				}
				cmd.expiryOption = store.ExpiryOptionExpiryUnixTimeInMilliseconds
				needTimeValue = true
			case "KEEPTTL":
				if cmd.expiryOption != store.ExpiryOptionNone {
					return nil, NewSyntaxError()
				}
				cmd.expiryOption = store.ExpiryOptionExpiryKeepTTL
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

type SetCommand struct {
	key             string
	value           string
	get             bool
	existenceOption existenceOption
	expiryOption    store.ExpiryOption
	expiry          int64
}

func (cmd SetCommand) Execute(s store.Store) (protocol.Data, error) {
	exists := s.Exists(cmd.key)
	var oldValue protocol.Data
	if exists && cmd.get {
		oldText, err := s.ReadString(cmd.key)
		if err != nil {
			return nil, err
		}
		oldValue = protocol.NewBulkString(oldText)
	}

	if exists && cmd.existenceOption == existenceOptionSetOnlyIfMissing {
		if cmd.get {
			return oldValue, nil
		}
		return nil, nil
	}
	if !exists && cmd.existenceOption == existenceOptionSetOnlyIfPresent {
		return nil, nil
	}

	s.Write(cmd.key, cmd.value)

	if cmd.get {
		return oldValue, nil
	}
	return protocol.NewSimpleString("OK"), nil
}
