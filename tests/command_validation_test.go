package tests_test

import (
	"redis-challenge/internal/command"
	"redis-challenge/internal/protocol"
	"redis-challenge/tests"
	"redis-challenge/tests/call"
	"testing"
)

type validationTestCases map[string]struct {
	calls        []call.DataCall
	driverChoice SelectTestCaseDriver
}

func TestCommandValidation(t *testing.T) {

	pingTestCases := validationTestCases{
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
	}

	echoTestCases := validationTestCases{
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
	}

	unknownCommandTestCases := validationTestCases{
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

	allTestCases := []validationTestCases{
		pingTestCases,
		echoTestCases,
		unknownCommandTestCases,
	}

	for _, testCases := range allTestCases {
		for name, testCase := range testCases {
			t.Run(name, func(t *testing.T) {
				validateCommands(t, testCase.calls, testCase.driverChoice)
			})
		}
	}
}

type SelectTestCaseDriver string

const (
	SelectTestCaseDriverRedisServer SelectTestCaseDriver = "redis-server-driver"
	SelectTestCaseDriverRedisClone  SelectTestCaseDriver = "redis-clone-driver"
)

func validateCommands(t testing.TB, calls []call.DataCall, driverChoice SelectTestCaseDriver) {
	switch driverChoice {
	case SelectTestCaseDriverRedisServer:
		tests.DriveProtocolAgainstServer(t, calls, tests.UseRealRedisServer)
	case SelectTestCaseDriverRedisClone:
		tests.DriveProtocolAgainstServer(t, calls, tests.UseChallengeServer)
	default:
		validateAgainstCommandValidator(t, calls)
	}
}

func validateAgainstCommandValidator(t testing.TB, calls []call.DataCall) {
	for _, c := range calls {
		_, errorData := command.Validate(c.RequestData())

		c.ConfirmValidation(t, errorData)
	}
}
