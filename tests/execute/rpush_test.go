package command_test

import (
	nanoid "github.com/matoous/go-nanoid/v2"
	"redis-challenge/internal/protocol"
	"redis-challenge/tests"
	"redis-challenge/tests/call"
	"testing"
)

func TestRightPush(t *testing.T) {

	uniqueSuffix := "-" + nanoid.Must(6)

	testCases := executionTestCases{
		"rpush 1 value to non-exists list creates the list with one value": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("RPUSH"),
						protocol.NewBulkString("rpush-key-to-empty" + uniqueSuffix),
						protocol.NewBulkString("1"),
					},
					protocol.NewSimpleInteger(1),
				),
			},
		},
		"rpush 2 values to non-exists list creates the list with both value": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("RPUSH"),
						protocol.NewBulkString("rpush-key-2-values" + uniqueSuffix),
						protocol.NewBulkString("1"),
						protocol.NewBulkString("2"),
					},
					protocol.NewSimpleInteger(2),
				),
			},
		},
		"rpush returns the total number of values in the list": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("RPUSH"),
						protocol.NewBulkString("rpush-key-returns-item-count" + uniqueSuffix),
						protocol.NewBulkString("a"),
						protocol.NewBulkString("b"),
					},
					protocol.NewSimpleInteger(2),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("RPUSH"),
						protocol.NewBulkString("rpush-key-returns-item-count" + uniqueSuffix),
						protocol.NewBulkString("c"),
						protocol.NewBulkString("d"),
						protocol.NewBulkString("e"),
					},
					protocol.NewSimpleInteger(5),
				),
			},
		},
		"rpush pushes to the end of the list in order": {
			calls: []call.DataCall{
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("RPUSH"),
						protocol.NewBulkString("rpush-key-to-end" + uniqueSuffix),
						protocol.NewBulkString("a"),
						protocol.NewBulkString("b"),
					},
				),
				call.NewFromDataWithoutError(
					[]protocol.Data{
						protocol.NewBulkString("RPUSH"),
						protocol.NewBulkString("rpush-key-to-end" + uniqueSuffix),
						protocol.NewBulkString("c"),
						protocol.NewBulkString("d"),
						protocol.NewBulkString("e"),
					},
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LRANGE"),
						protocol.NewBulkString("rpush-key-to-end" + uniqueSuffix),
						protocol.NewBulkString("0"),
						protocol.NewBulkString("-1"),
					},
					protocol.NewArray([]protocol.Data{
						protocol.NewBulkString("a"),
						protocol.NewBulkString("b"),
						protocol.NewBulkString("c"),
						protocol.NewBulkString("d"),
						protocol.NewBulkString("e"),
					}),
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
