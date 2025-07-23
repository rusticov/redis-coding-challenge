package validate_test

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/tests"
	"redis-challenge/tests/call"
	"testing"
)

func TestDelValidation(t *testing.T) {
	testCases := validationTestCases{
		"del command with no arguments has the wrong length": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("DEL"),
					},
					protocol.NewSimpleError("ERR wrong number of arguments for 'del' command"),
				),
			},
		},
		"del command with simple string key to delete has bad type": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("DEL"),
						protocol.NewSimpleString("key"),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got '+'"),
				),
			},
		},
		"del command with bulk string key to delete is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("DEL"),
						protocol.NewBulkString("key"),
					},
				),
			},
		},
		"del command with bulk string followed by simple string key to delete has bad type": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("DEL"),
						protocol.NewBulkString("key1"),
						protocol.NewSimpleString("key2"),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got '+'"),
				),
			},
		},
		"del command with sequence of only bulk strings is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("DEL"),
						protocol.NewBulkString("key1"),
						protocol.NewBulkString("key2"),
					},
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
