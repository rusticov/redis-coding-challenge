package validate_test

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/tests"
	"redis-challenge/tests/call"
	"testing"
)

func TestUnknownCommandValidation(t *testing.T) {
	testCases := map[string]struct {
		calls        []call.DataCall
		driverChoice tests.SelectTestCaseDriver
	}{
		"command 'UNKNOWN' is not a valid command": {
			calls: []call.DataCall{
				call.NewFromDataWithPartialError(
					[]protocol.Data{
						protocol.NewBulkString("UNKNOWN"),
					},
					"ERR unknown command 'UNKNOWN'",
				),
			},
		},
		"command 'BAD' is not a valid command": {
			calls: []call.DataCall{
				call.NewFromDataWithPartialError(
					[]protocol.Data{
						protocol.NewBulkString("BAD"),
					},
					"ERR unknown command 'BAD'",
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
