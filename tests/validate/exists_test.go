package validate_test

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/tests"
	"redis-challenge/tests/call"
	"testing"
)

func TestExistsValidation(t *testing.T) {
	testCases := validationTestCases{
		"exists command with no arguments has the wrong length": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("EXISTS"),
					},
					protocol.NewSimpleError("ERR wrong number of arguments for 'exists' command"),
				),
			},
		},
		"exists command with simple string key to exists has bad type": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("EXISTS"),
						protocol.NewSimpleString("key"),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got '+'"),
				),
			},
		},
		"exists command with bulk string key is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("EXISTS"),
						protocol.NewBulkString("key"),
					},
				),
			},
		},
		"exists command with bulk string followed by simple string key has bad type": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("EXISTS"),
						protocol.NewBulkString("key1"),
						protocol.NewSimpleString("key2"),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got '+'"),
				),
			},
		},
		"exists command with sequence of only bulk strings is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("EXISTS"),
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
