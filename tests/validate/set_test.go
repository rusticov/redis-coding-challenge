package validate_test

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/tests"
	"redis-challenge/tests/call"
	"testing"
)

func TestSetValidation(t *testing.T) {
	testCases := validationTestCases{
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

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			tests.ValidateCommands(t, testCase.calls, testCase.driverChoice)
		})
	}
}
