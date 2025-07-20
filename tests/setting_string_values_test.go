package tests_test

import (
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

	testCases := executionTestCases{
		"getting value that has been set (key-1)": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-1"),
						protocol.NewBulkString("value 1"),
					},
					protocol.NewSimpleString("OK"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("key-1"),
					},
					protocol.NewBulkString("value 1"),
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
