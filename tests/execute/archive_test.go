package command_test

import (
	"bytes"
	nanoid "github.com/matoous/go-nanoid/v2"
	"github.com/stretchr/testify/require"
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/server"
	"redis-challenge/internal/store"
	"redis-challenge/tests"
	"redis-challenge/tests/call"
	"testing"
	"time"
)

func TestWritingToArchive(t *testing.T) {

	uniqueSuffix := "-" + nanoid.Must(6)

	testCases := map[string]struct {
		calls              []call.DataCall
		delayBeforeRestore time.Duration
		postRestoreCalls   []call.DataCall
		driverChoice       tests.ServerVariant
	}{
		"getting value that has been set": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-with-value" + uniqueSuffix),
						protocol.NewBulkString("value 1"),
					},
					protocol.NewSimpleString("OK"),
				),
			},
			postRestoreCalls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("key-with-value" + uniqueSuffix),
					},
					protocol.NewBulkString("value 1"),
				),
			},
		},
		"getting expired value that has been set": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-with-expired-value" + uniqueSuffix),
						protocol.NewBulkString("value 1"),
						protocol.NewBulkString("EX"),
						protocol.NewBulkString("59"),
					},
					protocol.NewSimpleString("OK"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("key-with-expired-value" + uniqueSuffix),
					},
					nil,
				).WithDelay(time.Minute),
			},
			postRestoreCalls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("key-with-expired-value" + uniqueSuffix),
					},
					nil,
				),
			},
		},
		"getting value that has been set with relative expiry that was not expired during before restore, but is expired after restore given time has elapsed up to restore": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-with-relative-expiry" + uniqueSuffix),
						protocol.NewBulkString("value 1"),
						protocol.NewBulkString("EX"),
						protocol.NewBulkString("59"),
					},
					protocol.NewSimpleString("OK"),
				),
			},
			delayBeforeRestore: time.Minute,
			postRestoreCalls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("key-with-relative-expiry" + uniqueSuffix),
					},
					nil,
				),
			},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			// Given a log of executed commands is compiled
			clock := &store.FixedClock{TimeInMilliseconds: 1_000}

			buffer := bytes.NewBuffer(nil)
			tests.DriveProtocolAgainstServer(t, testCase.calls, testCase.driverChoice, buffer, clock)

			// When a server is restored after a delay
			clock.AddMilliseconds(testCase.delayBeforeRestore.Milliseconds())

			restoredServer, err := server.NewChallengeServer(0, store.NewBuilder().WithClock(clock)).
				RestoreFromArchive(buffer).
				Start()
			require.NoError(t, err)

			// Then postRestoreCalls can be successfully made
			tests.SendCallsToServer(t, restoredServer, testCase.postRestoreCalls, testCase.driverChoice, clock)
		})
	}
}
