package protocol_test

import (
	"github.com/stretchr/testify/assert"
	"redis-challenge/internal/protocol"
	"testing"
)

func TestDataTypeHasResp_Symbol(t *testing.T) {

	testCases := map[string]struct {
		data   protocol.Data
		symbol rune
	}{
		"simple string": {
			data:   protocol.NewSimpleString("message"),
			symbol: '+',
		},
		"simple error": {
			data:   protocol.NewSimpleError("message"),
			symbol: '-',
		},
		"simple integer": {
			data:   protocol.NewSimpleInteger(100),
			symbol: ':',
		},
		"bulk string": {
			data:   protocol.NewBulkString("message"),
			symbol: '$',
		},
		"array": {
			data:   protocol.NewArray(nil),
			symbol: '*',
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, protocol.DataTypeSymbol(testCase.symbol), testCase.data.Symbol(),
				"should return the correct symbol for the data type")
		})
	}
}
