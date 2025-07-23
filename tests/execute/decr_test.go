package command_test

import (
	nanoid "github.com/matoous/go-nanoid/v2"
	"redis-challenge/internal/protocol"
	"redis-challenge/tests"
	"redis-challenge/tests/call"
	"testing"
)

func TestDecrCommand(t *testing.T) {

	uniqueSuffix := "-" + nanoid.Must(6)

	testCases := executionTestCases{
		"decrementing a value that is an integer should decrement it in the store and return the new value": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-to-decrement" + uniqueSuffix),
						protocol.NewBulkString("5"),
					},
					protocol.NewSimpleString("OK"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("DECR"),
						protocol.NewBulkString("key-to-decrement" + uniqueSuffix),
					},
					protocol.NewSimpleInteger(4),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("key-to-decrement" + uniqueSuffix),
					},
					protocol.NewBulkString("4"),
				),
			},
		},
		"decrementing a key that does not exist should set value to 1 and return 1": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("DECR"),
						protocol.NewBulkString("unknown-key-to-decrement" + uniqueSuffix),
					},
					protocol.NewSimpleInteger(-1),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("unknown-key-to-decrement" + uniqueSuffix),
					},
					protocol.NewBulkString("-1"),
				),
			},
		},
		"decrementing a key that with non-integer value should error and not change the value": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-with-text-to-decrement" + uniqueSuffix),
						protocol.NewBulkString("text value"),
					},
					protocol.NewSimpleString("OK"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("DECR"),
						protocol.NewBulkString("key-with-text-to-decrement" + uniqueSuffix),
					},
					protocol.NewSimpleError("ERR value is not an integer or out of range"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("key-with-text-to-decrement" + uniqueSuffix),
					},
					protocol.NewBulkString("text value"),
				),
			},
		},
		"decrementing a key that with integer value that is large than int64 maximum should error and not change the value": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-with-bigint-to-decrement" + uniqueSuffix),
						protocol.NewBulkString("9223372036854775808"),
					},
					protocol.NewSimpleString("OK"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("DECR"),
						protocol.NewBulkString("key-with-bigint-to-decrement" + uniqueSuffix),
					},
					protocol.NewSimpleError("ERR value is not an integer or out of range"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("key-with-bigint-to-decrement" + uniqueSuffix),
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
