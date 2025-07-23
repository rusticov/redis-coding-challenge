package validate_test

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/tests"
	"redis-challenge/tests/call"
	"testing"
)

func TestRightPushValidation(t *testing.T) {
	testCases := validationTestCases{
		"rpush command with no arguments has the wrong length": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("RPUSH"),
					},
					protocol.NewSimpleError("ERR wrong number of arguments for 'rpush' command"),
				),
			},
		},
		"rpush command with only bulk string key has the wrong length": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("RPUSH"),
						protocol.NewBulkString("list-key"),
					},
					protocol.NewSimpleError("ERR wrong number of arguments for 'rpush' command"),
				),
			},
		},
		"rpush command with simple string key to delete has bad type": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("RPUSH"),
						protocol.NewSimpleString("list-key"),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got '+'"),
				),
			},
		},
		"rpush command with bulk string list argument is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("RPUSH"),
						protocol.NewBulkString("list-key"),
						protocol.NewBulkString("value"),
					},
				),
			},
		},
		"rpush command with simple string list argument has bad type": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("RPUSH"),
						protocol.NewBulkString("list-key"),
						protocol.NewSimpleString("value"),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got '+'"),
				),
			},
		},
		"rpush command with simple integer list argument is ok": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("RPUSH"),
						protocol.NewBulkString("list-key"),
						protocol.NewSimpleInteger(42),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got ':'"),
				),
			},
		},
		"rpush command with multiple mixed values should give bad type error": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("RPUSH"),
						protocol.NewBulkString("list-key"),
						protocol.NewBulkString("value"),
						protocol.NewSimpleInteger(42),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got ':'"),
				),
			},
		},
		"rpush command with multiple bulk strings is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("RPUSH"),
						protocol.NewBulkString("list-key"),
						protocol.NewBulkString("value 1"),
						protocol.NewBulkString("value 2"),
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
