package validate_test

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/tests"
	"redis-challenge/tests/call"
	"testing"
)

func TestLeftPushValidation(t *testing.T) {
	testCases := map[string]struct {
		calls        []call.DataCall
		driverChoice tests.SelectTestCaseDriver
	}{
		"lpush command with no arguments has the wrong length": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LPUSH"),
					},
					protocol.NewSimpleError("ERR wrong number of arguments for 'lpush' command"),
				),
			},
		},
		"lpush command with only bulk string key has the wrong length": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LPUSH"),
						protocol.NewBulkString("list-key"),
					},
					protocol.NewSimpleError("ERR wrong number of arguments for 'lpush' command"),
				),
			},
		},
		"lpush command with simple string key to delete has bad type": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LPUSH"),
						protocol.NewSimpleString("list-key"),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got '+'"),
				),
			},
		},
		"lpush command with bulk string list argument is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("LPUSH"),
						protocol.NewBulkString("list-key"),
						protocol.NewBulkString("value"),
					},
				),
			},
		},
		"lpush command with simple string list argument has bad type": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LPUSH"),
						protocol.NewBulkString("list-key"),
						protocol.NewSimpleString("value"),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got '+'"),
				),
			},
		},
		"lpush command with simple integer list argument is ok": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LPUSH"),
						protocol.NewBulkString("list-key"),
						protocol.NewSimpleInteger(42),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got ':'"),
				),
			},
		},
		"lpush command with multiple mixed values should give bad type error": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LPUSH"),
						protocol.NewBulkString("list-key"),
						protocol.NewBulkString("value"),
						protocol.NewSimpleInteger(42),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got ':'"),
				),
			},
		},
		"lpush command with multiple bulk strings is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("LPUSH"),
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
