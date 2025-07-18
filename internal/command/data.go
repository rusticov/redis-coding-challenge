package command

import (
	"fmt"
	"redis-challenge/internal/protocol"
)

type Data struct {
	Name      string
	Arguments []protocol.Data
}

func FromData(data protocol.Data) (Data, error) {
	if array, ok := data.(protocol.Array); ok {
		if len(array.Data) == 0 {
			return Data{}, fmt.Errorf("missing command name")
		}

		if name, ok := array.Data[0].(protocol.BulkString); ok {
			return Data{
				Name:      string(name),
				Arguments: array.Data[1:],
			}, nil
		}

		return Data{}, fmt.Errorf("command name must be a bulk string")
	}
	return Data{}, fmt.Errorf("not a command")
}
