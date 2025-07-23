package command_test

import (
	nanoid "github.com/matoous/go-nanoid/v2"
	"redis-challenge/internal/protocol"
	"redis-challenge/tests"
	"redis-challenge/tests/call"
	"testing"
)

func TestLeftRange(t *testing.T) {

	uniqueSuffix := "-" + nanoid.Must(6)

	testCases := executionTestCases{
		"lrange 0 0 returns value added last to the left": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LPUSH"),
						protocol.NewBulkString("key-for-single-value" + uniqueSuffix),
						protocol.NewBulkString("one"),
						protocol.NewBulkString("two"),
					},
					protocol.NewSimpleInteger(2),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LRANGE"),
						protocol.NewBulkString("key-for-single-value" + uniqueSuffix),
						protocol.NewBulkString("0"),
						protocol.NewBulkString("0"),
					},
					protocol.NewArray([]protocol.Data{
						protocol.NewBulkString("two"),
					}),
				),
			},
		},
		"lrange 1 3 returns range from the middle in reversed order (as values added in reversed order)": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LPUSH"),
						protocol.NewBulkString("key-for-mid-range" + uniqueSuffix),
						protocol.NewBulkString("a"),
						protocol.NewBulkString("b"),
						protocol.NewBulkString("c"),
						protocol.NewBulkString("d"),
						protocol.NewBulkString("e"),
					},
					protocol.NewSimpleInteger(5),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LRANGE"),
						protocol.NewBulkString("key-for-mid-range" + uniqueSuffix),
						protocol.NewBulkString("1"),
						protocol.NewBulkString("3"),
					},
					protocol.NewArray([]protocol.Data{
						protocol.NewBulkString("d"),
						protocol.NewBulkString("c"),
						protocol.NewBulkString("b"),
					}),
				),
			},
		},
	}

	for name, testCase := range testCases {
		t.Skip("awaiting a new implementation of LRANGE command")
		t.Run(name, func(t *testing.T) {
			tests.DriveProtocolAgainstServer(t, testCase.calls, testCase.driverChoice)
		})
	}
}
