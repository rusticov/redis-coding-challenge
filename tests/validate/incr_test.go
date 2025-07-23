package validate_test

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/tests"
	"redis-challenge/tests/call"
	"testing"
)

func TestIncrValidation(t *testing.T) {
	testCases := validationTestCases{
		"incr command with no arguments has the wrong length": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("INCR"),
					},
					protocol.NewSimpleError("ERR wrong number of arguments for 'incr' command"),
				),
			},
		},
		"incr command with simple string key to increment has bad type": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("INCR"),
						protocol.NewSimpleString("key"),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got '+'"),
				),
			},
		},
		"incr command with bulk string key to increment is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("INCR"),
						protocol.NewBulkString("key"),
					},
				),
			},
		},
		"incr command with bulk string key and integer has the wrong type": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("INCR"),
						protocol.NewBulkString("key"),
						protocol.NewSimpleInteger(1),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got ':'"),
				),
			},
		},
		"incr command with 2 bulk string key to increment has wrong length": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("INCR"),
						protocol.NewBulkString("key1"),
						protocol.NewBulkString("key2"),
					},
					protocol.NewSimpleError("ERR wrong number of arguments for 'incr' command"),
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
