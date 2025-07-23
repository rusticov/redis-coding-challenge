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

	getCommandTestCases := validationTestCases{
		"get command with no arguments has the wrong length": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
					},
					protocol.NewSimpleError("ERR wrong number of arguments for 'get' command"),
				),
			},
		},
		"get command with bulk string message is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("key"),
					},
				),
			},
		},
		"get command with multiple arguments has the wrong length": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("message"),
						protocol.NewBulkString("message"),
					},
					protocol.NewSimpleError("ERR wrong number of arguments for 'get' command"),
				),
			},
		},
		"get command with simple string message": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
						protocol.NewSimpleString("message"),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got '+'"),
				),
			},
		},
	}

	delCommandTestCases := validationTestCases{
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

	setCommandTestCases := validationTestCases{
		"set command with no arguments has the wrong length": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
					},
					protocol.NewSimpleError("ERR wrong number of arguments for 'set' command"),
				),
			},
		},
		"set command with only key has the wrong length": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
					},
					protocol.NewSimpleError("ERR wrong number of arguments for 'set' command"),
				),
			},
		},
		"set command with only key as a simple string has the wrong length": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewSimpleString("key"),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got '+'"),
				),
			},
		},
		"set command with only key and bulk string value is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("value"),
					},
				),
			},
		},
		"set command with only key and an array of bulk strings value is an error": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewArray(
							[]protocol.Data{protocol.NewBulkString("value")},
						),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got '*'"),
				),
			},
		},
		"set command with only key and simple string value is an error": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewSimpleString("value"),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got '+'"),
				),
			},
		},
		"set command with only key and simple integer value is and error": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewSimpleInteger(42),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got ':'"),
				),
			},
		},
		"set command with only key and simple error value is an error": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewSimpleError("error value"),
					},
					protocol.NewSimpleError("ERR Protocol error: expected '$', got '-'"),
				),
			},
		},
		"set command with only key and bulk string value and GET is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("value"),
						protocol.NewBulkString("GET"),
					},
				),
			},
		},
		"set command with only key and bulk string value and NX is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("value"),
						protocol.NewBulkString("NX"),
					},
				),
			},
		},
		"set command with only key and bulk string value and XX is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("value"),
						protocol.NewBulkString("XX"),
					},
				),
			},
		},
		"set command with only key and bulk string value and GET and NX is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("value"),
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("NX"),
					},
				),
			},
		},
		"set command with only key and bulk string value and NX and GET is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("value"),
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("NX"),
					},
				),
			},
		},
		"set command with only key and bulk string value and XX and NX is an error": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("value"),
						protocol.NewBulkString("XX"),
						protocol.NewBulkString("NX"),
					},
					protocol.NewSimpleError("ERR syntax error"),
				),
			},
		},
		"set command with only key and bulk string value and NX and XX is an error": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("value"),
						protocol.NewBulkString("NX"),
						protocol.NewBulkString("XX"),
					},
					protocol.NewSimpleError("ERR syntax error"),
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
		delCommandTestCases,
		getCommandTestCases,
		setCommandTestCases,
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
