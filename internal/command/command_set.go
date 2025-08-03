package command

import (
	"bytes"
	"log/slog"
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
	"strconv"
)

type SetValidator struct {
	clock store.Clock
}

func (v *SetValidator) Validate(requestBytes []byte, arguments []protocol.Data) (Command, protocol.Data) {
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
		requestBytes: requestBytes,
		key:          string(arguments[0].(protocol.BulkString)),
		value:        string(arguments[1].(protocol.BulkString)),
		expiryOption: store.ExpiryOptionNone,
		expiry:       0,
	}

	var expiryIsRelative bool
	var expiryOptionIndex int
	var needTimeValue bool

	for i, arg := range arguments[2:] {
		if bulkText, ok := arg.(protocol.BulkString); ok {
			if needTimeValue {
				expiry, err := strconv.ParseInt(string(bulkText), 10, 64)
				if err != nil {
					return nil, NewSyntaxError()
				}

				needTimeValue = false
				cmd.expiryOption, cmd.expiry = ExpiryTimestamp(v.clock, cmd.expiryOption, expiry)
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

				expiryIsRelative = true
				expiryOptionIndex = i + 2
			case "PX":
				if cmd.expiryOption != store.ExpiryOptionNone {
					return nil, NewSyntaxError()
				}
				cmd.expiryOption = store.ExpiryOptionExpiryMilliseconds
				needTimeValue = true

				expiryIsRelative = true
				expiryOptionIndex = i + 2
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

	if expiryIsRelative {
		if updatedRequestBytes, ok := v.requestBytesWithUpdatedExpiry(arguments, expiryOptionIndex, cmd.expiry); ok {
			cmd.requestBytes = updatedRequestBytes
		}
	}

	return cmd, nil
}

func (v *SetValidator) requestBytesWithUpdatedExpiry(arguments []protocol.Data, expiryOptionIndex int, expiry int64) ([]byte, bool) {
	newArguments := make([]protocol.Data, len(arguments)+1)
	copy(newArguments[1:], arguments)

	newArguments[0] = protocol.NewBulkString("SET")
	newArguments[expiryOptionIndex+1] = protocol.NewBulkString("PXAT")
	newArguments[expiryOptionIndex+2] = protocol.NewBulkString(strconv.FormatInt(expiry, 10))

	buffer := bytes.NewBuffer(nil)
	err := protocol.WriteData(buffer, protocol.Array{Data: newArguments})

	if err != nil {
		slog.Error("cannot convert the request to convert expiry", "error", err)
		return nil, false
	} else {
		return buffer.Bytes(), true
	}
}

type existenceOption string

const (
	existenceOptionSetOnlyIfMissing existenceOption = "NX"
	existenceOptionSetOnlyIfPresent existenceOption = "XX"
	existenceOptionNone             existenceOption = ""
)

type SetCommand struct {
	requestBytes    []byte
	key             string
	value           string
	get             bool
	existenceOption existenceOption
	expiryOption    store.ExpiryOption
	expiry          int64
}

func (cmd SetCommand) Request() ([]byte, Type) {
	return cmd.requestBytes, TypeUpdate
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

	s.Write(cmd.key, cmd.value, cmd.expiryOption, cmd.expiry)

	if cmd.get {
		return oldValue, nil
	}
	return protocol.NewSimpleString("OK"), nil
}
