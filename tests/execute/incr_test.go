package command_test

import (
	nanoid "github.com/matoous/go-nanoid/v2"
	"redis-challenge/internal/protocol"
	"redis-challenge/tests"
	"redis-challenge/tests/call"
	"testing"
)

func TestIncrCommand(t *testing.T) {

	uniqueSuffix := "-" + nanoid.Must(6)

	testCases := executionTestCases{
		"incrementing a value that is an integer should increment it in the store and return the new value": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-to-increment" + uniqueSuffix),
						protocol.NewBulkString("5"),
					},
					protocol.NewSimpleString("OK"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("INCR"),
						protocol.NewBulkString("key-to-increment" + uniqueSuffix),
					},
					protocol.NewSimpleInteger(6),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("key-to-increment" + uniqueSuffix),
					},
					protocol.NewBulkString("6"),
				),
			},
		},
		"incrementing a key that does not exist should set value to 1 and return 1": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("INCR"),
						protocol.NewBulkString("unknown-key-to-increment" + uniqueSuffix),
					},
					protocol.NewSimpleInteger(1),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("unknown-key-to-increment" + uniqueSuffix),
					},
					protocol.NewBulkString("1"),
				),
			},
		},
		"incrementing a key that with non-integer value should error and not change the value": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-with-text-to-increment" + uniqueSuffix),
						protocol.NewBulkString("text value"),
					},
					protocol.NewSimpleString("OK"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("INCR"),
						protocol.NewBulkString("key-with-text-to-increment" + uniqueSuffix),
					},
					protocol.NewSimpleError("ERR value is not an integer or out of range"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("key-with-text-to-increment" + uniqueSuffix),
					},
					protocol.NewBulkString("text value"),
				),
			},
		},
		"incrementing a key that with integer value that is large than int64 maximum should error and not change the value": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-with-bigint-to-increment" + uniqueSuffix),
						protocol.NewBulkString("9223372036854775808"),
					},
					protocol.NewSimpleString("OK"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("INCR"),
						protocol.NewBulkString("key-with-bigint-to-increment" + uniqueSuffix),
					},
					protocol.NewSimpleError("ERR value is not an integer or out of range"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("key-with-bigint-to-increment" + uniqueSuffix),
					},
					protocol.NewBulkString("9223372036854775808"),
				),
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			tests.DriveProtocolAgainstServer(t, testCase.calls, testCase.driverChoice)
		})
	}
}
