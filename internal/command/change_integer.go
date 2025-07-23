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
	valueFromStore, exists := s.Get(cmd.key)

	if !exists {
		s.LoadOrStore(cmd.key, fmt.Sprintf("%d", cmd.change)) // TODO test failure to set here
		return protocol.NewSimpleInteger(cmd.change), nil
	}

	value, ok := parseStoreValueAsInt(valueFromStore)
	if !ok {
		return protocol.NewSimpleError("ERR value is not an integer or out of range"), nil
	}
	value += cmd.change

	s.CompareAndSwap(cmd.key, valueFromStore, strconv.FormatInt(value, 10)) // TODO test failure to set here

	return protocol.NewSimpleInteger(value), nil
}

func parseStoreValueAsInt(data any) (int64, bool) {
	if text, ok := data.(string); ok {
		value, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return 0, false
		}
		return value, true
	}
	return 0, false
}
