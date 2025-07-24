package command

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
)

type ChangeIntegerCommand struct {
	key    string
	change int64
}

func (cmd ChangeIntegerCommand) Execute(s store.Store) (protocol.Data, error) {
	value, err := s.Increment(cmd.key, cmd.change)
	if err != nil {
		return protocol.NewSimpleError("ERR value is not an integer or out of range"), nil
	}

	return protocol.NewSimpleInteger(value), nil
}
