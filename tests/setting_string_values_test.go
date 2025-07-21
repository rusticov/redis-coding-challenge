package tests_test

import (
	nanoid "github.com/matoous/go-nanoid/v2"
	"redis-challenge/internal/protocol"
	"redis-challenge/tests"
	"redis-challenge/tests/call"
	"testing"
)

type executionTestCases map[string]struct {
	calls        []call.DataCall
	driverChoice tests.ServerVariant
}

func TestSettingStringValues(t *testing.T) {

	uniqueSuffix := "-" + nanoid.Must(6)

	testCases := executionTestCases{
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
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("key-with-value" + uniqueSuffix),
					},
					protocol.NewBulkString("value 1"),
				),
			},
		},
		"getting value that has not been set": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("key-no-value" + uniqueSuffix),
					},
					nil,
				),
			},
		},
		"setting value with the get option returns the previous value": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-with-value-get" + uniqueSuffix),
						protocol.NewBulkString("first value"),
						protocol.NewBulkString("GET"),
					},
					nil,
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-with-value-get" + uniqueSuffix),
						protocol.NewBulkString("second value"),
						protocol.NewBulkString("GET"),
					},
					protocol.NewBulkString("first value"),
				),
			},
		},
		"setting value with NX options only sets the value if it does not exist": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-with-value-nx" + uniqueSuffix),
						protocol.NewBulkString("first value"),
						protocol.NewBulkString("NX"),
					},
					protocol.NewSimpleString("OK"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-with-value-nx" + uniqueSuffix),
						protocol.NewBulkString("second value"),
						protocol.NewBulkString("NX"),
					},
					nil,
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("key-with-value-nx" + uniqueSuffix),
					},
					protocol.NewBulkString("first value"),
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
