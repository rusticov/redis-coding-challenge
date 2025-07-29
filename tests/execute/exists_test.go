package command_test

import (
	nanoid "github.com/matoous/go-nanoid/v2"
	"redis-challenge/internal/protocol"
	"redis-challenge/tests"
	"redis-challenge/tests/call"
	"testing"
)

func TestCheckingForExistence(t *testing.T) {

	uniqueSuffix := "-" + nanoid.Must(6)

	testCases := map[string]struct {
		calls        []call.DataCall
		driverChoice tests.ServerVariant
	}{
		"exists for a key that is unknown should return 0": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("EXISTS"),
						protocol.NewBulkString("random-key-with-no-value" + uniqueSuffix),
					},
					protocol.NewSimpleInteger(0),
				),
			},
		},
		"exists for a known key should return 1 and leave key unchanged in the store": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-to-check" + uniqueSuffix),
						protocol.NewBulkString("value 1"),
					},
					protocol.NewSimpleString("OK"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("EXISTS"),
						protocol.NewBulkString("key-to-check" + uniqueSuffix),
					},
					protocol.NewSimpleInteger(1),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("key-to-check" + uniqueSuffix),
					},
					protocol.NewBulkString("value 1"),
				),
			},
		},
		"exists for a mix of known and unknown keys should return count known keys": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-to-check-1" + uniqueSuffix),
						protocol.NewBulkString("value 1"),
					},
					protocol.NewSimpleString("OK"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-to-check-3" + uniqueSuffix),
						protocol.NewBulkString("value 3"),
					},
					protocol.NewSimpleString("OK"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("DEL"),
						protocol.NewBulkString("key-to-check-1" + uniqueSuffix),
						protocol.NewBulkString("key-to-check-2" + uniqueSuffix),
						protocol.NewBulkString("key-to-check-3" + uniqueSuffix),
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
