package validate_test

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/tests"
	"redis-challenge/tests/call"
	"testing"
)

func TestPingValidation(t *testing.T) {
	testCases := validationTestCases{
		"ping command with no arguments is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("PING"),
					},
				),
			},
		},
		"ping command with message is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("PING"),
						protocol.NewBulkString("message"),
					},
				),
			},
		},
		"ping command with two message is too long": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("PING"),
						protocol.NewBulkString("message"),
						protocol.NewBulkString("message"),
					},
					protocol.NewSimpleError("ERR wrong number of arguments for 'ping' command"),
				),
			},
		},
		"ping command with simple string message": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("PING"),
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
