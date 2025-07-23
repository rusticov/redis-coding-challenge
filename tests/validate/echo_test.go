package validate_test

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/tests"
	"redis-challenge/tests/call"
	"testing"
)

type validationTestCases map[string]struct {
	calls        []call.DataCall
	driverChoice tests.SelectTestCaseDriver
}

func TestEchoValidation(t *testing.T) {
	testCases := validationTestCases{
		"echo command with no arguments has the wrong length": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("ECHO"),
					},
					protocol.NewSimpleError("ERR wrong number of arguments for 'echo' command"),
				),
			},
		},
		"echo command with bulk string message is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("ECHO"),
						protocol.NewBulkString("message"),
					},
				),
			},
		},
		"echo command with multiple arguments has the wrong length": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("ECHO"),
						protocol.NewBulkString("message"),
						protocol.NewBulkString("message"),
					},
					protocol.NewSimpleError("ERR wrong number of arguments for 'echo' command"),
				),
			},
		},
		"echo command with simple string message": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("ECHO"),
						protocol.NewSimpleString("message"),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got '+'"),
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
