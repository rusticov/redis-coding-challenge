package command_test

import (
	nanoid "github.com/matoous/go-nanoid/v2"
	"redis-challenge/internal/protocol"
	"redis-challenge/tests"
	"redis-challenge/tests/call"
	"testing"
	"time"
)

func TestExpiry(t *testing.T) {

	uniqueSuffix := "-" + nanoid.Must(6)

	testCases := map[string]struct {
		calls        []call.DataCall
		driverChoice tests.ServerVariant
	}{
		"getting value that has expired should return nil": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-expired" + uniqueSuffix),
						protocol.NewBulkString("value 1"),
						protocol.NewBulkString("EX"),
						protocol.NewBulkString("1"),
					},
					protocol.NewSimpleString("OK"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("key-expired" + uniqueSuffix),
					},
					nil,
				).WithDelay(time.Second + time.Millisecond),
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			tests.DriveProtocolAgainstServer(t, testCase.calls, testCase.driverChoice)
		})
	}
}
