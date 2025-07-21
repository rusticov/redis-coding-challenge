package command

import (
	"fmt"
	"redis-challenge/internal/protocol"
)

func NewWrongDataTypeError(data protocol.Data, expectedType protocol.DataTypeSymbol) protocol.SimpleError {
	return protocol.NewSimpleError(fmt.Sprintf("ERR Protocol error: expected '%c', got '%c'", expectedType, data.Symbol()))
}

func NewSyntaxError() protocol.SimpleError {
	return protocol.NewSimpleError(fmt.Sprintf("ERR syntax error"))
}
