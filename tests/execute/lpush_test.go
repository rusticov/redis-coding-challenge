package command_test

import (
	nanoid "github.com/matoous/go-nanoid/v2"
	"redis-challenge/internal/protocol"
	"redis-challenge/tests"
	"redis-challenge/tests/call"
	"testing"
)

func TestLeftPush(t *testing.T) {

	uniqueSuffix := "-" + nanoid.Must(6)

	testCases := executionTestCases{
		"lpush 1 value to non-exists list creates the list with one value": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LPUSH"),
						protocol.NewBulkString("key-to-empty" + uniqueSuffix),
						protocol.NewBulkString("1"),
					},
					protocol.NewSimpleInteger(1),
				),
			},
		},
		"lpush 2 values to non-exists list creates the list with both value": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LPUSH"),
						protocol.NewBulkString("key-to-empty" + uniqueSuffix),
						protocol.NewBulkString("1"),
						protocol.NewBulkString("2"),
					},
					protocol.NewSimpleInteger(2),
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
