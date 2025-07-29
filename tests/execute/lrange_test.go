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

	testCases := map[string]struct {
		calls        []call.DataCall
		driverChoice tests.ServerVariant
	}{
		"lrange 0 0 of a string value should be error": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-for-string" + uniqueSuffix),
						protocol.NewBulkString("text"),
					},
					protocol.NewSimpleString("OK"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LRANGE"),
						protocol.NewBulkString("key-for-string" + uniqueSuffix),
						protocol.NewBulkString("0"),
						protocol.NewBulkString("0"),
					},
					protocol.NewSimpleError("WRONGTYPE Operation against a key holding the wrong kind of value"),
				),
			},
		},
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
		"lrange with negative end counts -1 as the end and so -2 as one in from the end": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LPUSH"),
						protocol.NewBulkString("key-for-negative-end" + uniqueSuffix),
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
						protocol.NewBulkString("key-for-negative-end" + uniqueSuffix),
						protocol.NewBulkString("1"),
						protocol.NewBulkString("-2"),
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
		t.Run(name, func(t *testing.T) {
			tests.DriveProtocolAgainstServer(t, testCase.calls, testCase.driverChoice)
		})
	}
}
