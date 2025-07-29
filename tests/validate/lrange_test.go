package validate_test

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/tests"
	"redis-challenge/tests/call"
	"testing"
)

func TestLeftRangeValidation(t *testing.T) {
	testCases := map[string]struct {
		calls        []call.DataCall
		driverChoice tests.SelectTestCaseDriver
	}{
		"lrange command with no arguments has the wrong length": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LRANGE"),
					},
					protocol.NewSimpleError("ERR wrong number of arguments for 'lrange' command"),
				),
			},
		},
		"lrange command with simple string key has bad type": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LRANGE"),
						protocol.NewSimpleString("key"),
						protocol.NewBulkString("0"),
						protocol.NewBulkString("2"),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got '+'"),
				),
			},
		},
		"lrange command with integer first range value has bad type": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LRANGE"),
						protocol.NewBulkString("key"),
						protocol.NewSimpleInteger(0),
						protocol.NewBulkString("2"),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got ':'"),
				),
			},
		},
		"lrange command with integer second range value has bad type": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LRANGE"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("0"),
						protocol.NewSimpleInteger(2),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got ':'"),
				),
			},
		},
		"lrange command with bulk string and bulk strings for the range values is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("LRANGE"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("0"),
						protocol.NewBulkString("2"),
					},
				),
			},
		},
		"lrange command with non-integer left range is out of range": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LRANGE"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("zero"),
						protocol.NewBulkString("2"),
					},
					protocol.NewSimpleError("ERR value is not an integer or out of range"),
				),
			},
		},
		"lrange command with non-integer right range is out of range": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LRANGE"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("0"),
						protocol.NewBulkString("two"),
					},
					protocol.NewSimpleError("ERR value is not an integer or out of range"),
				),
			},
		},
		"lrange command with bulk string and 2 integers and extra value is too long": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LRANGE"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("0"),
						protocol.NewBulkString("2"),
						protocol.NewBulkString("4"),
					},
					protocol.NewSimpleError("ERR wrong number of arguments for 'lrange' command"),
				),
			},
		},
		"lrange command without stop value is too short": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LRANGE"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("0"),
					},
					protocol.NewSimpleError("ERR wrong number of arguments for 'lrange' command"),
				),
			},
		},
		"lrange command with only key and no range values is too short": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LRANGE"),
						protocol.NewBulkString("key"),
					},
					protocol.NewSimpleError("ERR wrong number of arguments for 'lrange' command"),
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
