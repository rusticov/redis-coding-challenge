package validate_test

import (
	"redis-challenge/internal/protocol"
	"redis-challenge/tests"
	"redis-challenge/tests/call"
	"testing"
)

func TestSetValidation(t *testing.T) {
	testCases := map[string]struct {
		calls        []call.DataCall
		driverChoice tests.SelectTestCaseDriver
	}{
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
		"set command with EX and no seconds argument is an error": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("value"),
						protocol.NewBulkString("EX"),
					},
					protocol.NewSimpleError("ERR syntax error"),
				),
			},
		},
		"set command with EX and seconds argument is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("value"),
						protocol.NewBulkString("EX"),
						protocol.NewBulkString("42"),
					},
				),
			},
		},
		"set command with EX and seconds argument followed by GET is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("value"),
						protocol.NewBulkString("EX"),
						protocol.NewBulkString("42"),
						protocol.NewBulkString("GET"),
					},
				),
			},
		},
		"set command with PX and no milliseconds argument is an error": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("value"),
						protocol.NewBulkString("PX"),
					},
					protocol.NewSimpleError("ERR syntax error"),
				),
			},
		},
		"set command with PX and milliseconds argument is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("value"),
						protocol.NewBulkString("PX"),
						protocol.NewBulkString("42"),
					},
				),
			},
		},
		"set command with mixing EX with seconds arguments followed by PX is an error": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("value"),
						protocol.NewBulkString("EX"),
						protocol.NewBulkString("42"),
						protocol.NewBulkString("PX"),
						protocol.NewBulkString("42000"),
					},
					protocol.NewSimpleError("ERR syntax error"),
				),
			},
		},
		"set command with mixing PX with milliseconds arguments followed by EX is an error": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("value"),
						protocol.NewBulkString("PX"),
						protocol.NewBulkString("42000"),
						protocol.NewBulkString("EX"),
						protocol.NewBulkString("42"),
					},
					protocol.NewSimpleError("ERR syntax error"),
				),
			},
		},
		"set command with EXAT and no seconds argument is an error": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("value"),
						protocol.NewBulkString("EXAT"),
					},
					protocol.NewSimpleError("ERR syntax error"),
				),
			},
		},
		"set command with EXAT and seconds argument is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("value"),
						protocol.NewBulkString("EXAT"),
						protocol.NewBulkString("42345345"),
					},
				),
			},
		},
		"set command with mixing EX with seconds arguments followed by EXAT is an error": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("value"),
						protocol.NewBulkString("EX"),
						protocol.NewBulkString("42"),
						protocol.NewBulkString("EXAT"),
						protocol.NewBulkString("42000"),
					},
					protocol.NewSimpleError("ERR syntax error"),
				),
			},
		},
		"set command with PXAT and no milliseconds argument is an error": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("value"),
						protocol.NewBulkString("EXAT"),
					},
					protocol.NewSimpleError("ERR syntax error"),
				),
			},
		},
		"set command with PXAT and milliseconds argument is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("value"),
						protocol.NewBulkString("PXAT"),
						protocol.NewBulkString("42345345"),
					},
				),
			},
		},
		"set command with mixing EX with seconds arguments followed by PXAT is an error": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("value"),
						protocol.NewBulkString("EX"),
						protocol.NewBulkString("42"),
						protocol.NewBulkString("PXAT"),
						protocol.NewBulkString("42000"),
					},
					protocol.NewSimpleError("ERR syntax error"),
				),
			},
		},
		"set command with KEEPTTL is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("value"),
						protocol.NewBulkString("KEEPTTL"),
					},
				),
			},
		},
		"set command with mixing EX with seconds arguments followed by KEEPTTL is an error": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("value"),
						protocol.NewBulkString("EX"),
						protocol.NewBulkString("42"),
						protocol.NewBulkString("KEEPTTL"),
					},
					protocol.NewSimpleError("ERR syntax error"),
				),
			},
		},
		"set command with mixing KEEPTTL followed by EX with seconds arguments is an error": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("value"),
						protocol.NewBulkString("KEEPTTL"),
						protocol.NewBulkString("EX"),
						protocol.NewBulkString("42"),
					},
					protocol.NewSimpleError("ERR syntax error"),
				),
			},
		},
		"set command with KEEPTTL followed by GET is ok": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key"),
						protocol.NewBulkString("value"),
						protocol.NewBulkString("KEEPTTL"),
						protocol.NewBulkString("GET"),
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
