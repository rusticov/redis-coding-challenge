package tests_test

import (
	"github.com/stretchr/testify/assert"
	"redis-challenge/internal/command"
	"redis-challenge/internal/protocol"
	"testing"
)

type validationTestCases map[string]struct {
	protocol      protocol.Data
	expectedError protocol.Data
	isOK          bool
}

func TestCommandValidation(t *testing.T) {

	pingTestCases := validationTestCases{
		"ping command with no arguments is ok": {
			protocol: protocol.Array{Data: []protocol.Data{
				protocol.NewBulkString("PING"),
			}},
			isOK: true,
		},
		"ping command with message is ok": {
			protocol: protocol.Array{Data: []protocol.Data{
				protocol.NewBulkString("PING"),
				protocol.NewBulkString("message"),
			}},
			isOK: true,
		},
		"ping command with two message is too long": {
			protocol: protocol.Array{Data: []protocol.Data{
				protocol.NewBulkString("PING"),
				protocol.NewBulkString("message"),
				protocol.NewBulkString("message"),
			}},
			expectedError: protocol.NewSimpleError("ERR wrong number of arguments for 'ping' command"),
		},
	}

	echoTestCases := validationTestCases{
		"echo command with no arguments has the wrong length": {
			protocol: protocol.Array{Data: []protocol.Data{
				protocol.NewBulkString("ECHO"),
			}},
			expectedError: protocol.NewSimpleError("ERR wrong number of arguments for 'echo' command"),
		},
		"echo command with bulk string message is ok": {
			protocol: protocol.Array{Data: []protocol.Data{
				protocol.NewBulkString("ECHO"),
				protocol.NewBulkString("message"),
			}},
			isOK: true,
		},
		"echo command with multiple arguments has the wrong length": {
			protocol: protocol.Array{Data: []protocol.Data{
				protocol.NewBulkString("ECHO"),
				protocol.NewBulkString("message"),
				protocol.NewBulkString("message"),
			}},
			expectedError: protocol.NewSimpleError("ERR wrong number of arguments for 'echo' command"),
		},
	}

	unknownCommandTestCases := validationTestCases{
		"command 'UNKNOWN' is not a valid command": {
			protocol: protocol.Array{Data: []protocol.Data{
				protocol.NewBulkString("UNKNOWN"),
			}},
			expectedError: protocol.NewSimpleError("ERR unknown command 'UNKNOWN'"),
		},
		"command 'BAD' is not a valid command": {
			protocol: protocol.Array{Data: []protocol.Data{
				protocol.NewBulkString("BAD"),
			}},
			expectedError: protocol.NewSimpleError("ERR unknown command 'BAD'"),
		},
	}

	allTestCases := []validationTestCases{
		pingTestCases,
		echoTestCases,
		unknownCommandTestCases,
	}

	for _, testCases := range allTestCases {
		for name, testCase := range testCases {
			t.Run(name, func(t *testing.T) {
				validateAgainstCommandValidator(t, testCase.protocol, testCase.expectedError)
			})
		}
	}
}

func validateAgainstCommandValidator(t testing.TB, input protocol.Data, expectedError protocol.Data) {
	_, errorData := command.Validate(input)

	if expectedError == nil {
		assert.Nil(t, errorData, "command should be valid")
	} else {
		assert.Equal(t, expectedError, errorData, "command should be invalid")
	}
}
