package validate_test

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/tests"
	"redis-challenge/tests/call"
	"testing"
)

func TestDecrValidation(t *testing.T) {
	testCases := validationTestCases{
		"decr command with no arguments has the wrong length": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("DECR"),
					},
					protocol.NewSimpleError("ERR wrong number of arguments for 'decr' command"),
				),
			},
		},
		"decr command with simple string key to decrement has bad type": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("DECR"),
						protocol.NewSimpleString("key"),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got '+'"),
				),
			},
		},
		"decr command with bulk string key to decrement is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("DECR"),
						protocol.NewBulkString("key"),
					},
				),
			},
		},
		"decr command with bulk string key and integer has the wrong type": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("DECR"),
						protocol.NewBulkString("key"),
						protocol.NewSimpleInteger(1),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got ':'"),
				),
			},
		},
		"decr command with 2 bulk string key to decrement has wrong length": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("DECR"),
						protocol.NewBulkString("key1"),
						protocol.NewBulkString("key2"),
					},
					protocol.NewSimpleError("ERR wrong number of arguments for 'decr' command"),
				),
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			tests.ValidateCommands(t, testCase.calls, testCase.driverChoice)
		})
	}
}
