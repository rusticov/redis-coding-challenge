package command

import (
	"redis-challenge/internal/protocol"
)

type Data struct {
	Name      string
	Arguments []protocol.Data
}

func FromData(data protocol.Data) (Data, protocol.Data) {
	if array, ok := data.(protocol.Array); ok {
		if len(array.Data) == 0 {
			return Data{}, protocol.NewSimpleError("missing command name")
		}

		nameData := array.Data[0]

		if name, ok := nameData.(protocol.BulkString); ok {
			return Data{
				Name:      string(name),
				Arguments: array.Data[1:],
			}, nil
		}

		return Data{}, NewWrongDataTypeError(nameData, protocol.BulkStringSymbol)
	}

	if simpleError, ok := data.(protocol.SimpleError); ok {
		return Data{}, simpleError
	}

	return Data{}, protocol.NewSimpleError("not a command")
}
