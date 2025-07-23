package command

import (
	"fmt"
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
	"strconv"
)

type ChangeIntegerCommand struct {
	key    string
	change int64
}

func (cmd ChangeIntegerCommand) Execute(s store.Store) (protocol.Data, error) {
	textValue, exists := s.Get(cmd.key)

	if !exists {
		s.LoadOrStore(cmd.key, fmt.Sprintf("%d", cmd.change)) // TODO test failure to set here
		return protocol.NewSimpleInteger(cmd.change), nil
	}

	var value int64
	if textValue != "" {
		var err error
		value, err = strconv.ParseInt(textValue, 10, 64)
		if err != nil {
			return protocol.NewSimpleError("ERR value is not an integer or out of range"), nil
		}
	}
	value += cmd.change

	s.CompareAndSwap(cmd.key, textValue, strconv.FormatInt(value, 10)) // TODO test failure to set here

	return protocol.NewSimpleInteger(value), nil
}
