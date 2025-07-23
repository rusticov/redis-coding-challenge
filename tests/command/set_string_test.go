package command_test

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
		"setting value with NX option only sets the value if it does not exist": {
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
		"setting value with GET and NX options get value when NX option can set the value": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-with-value-nx" + uniqueSuffix),
						protocol.NewBulkString("first value"),
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("NX"),
					},
					nil,
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-with-value-nx" + uniqueSuffix),
						protocol.NewBulkString("second value"),
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("NX"),
					},
					protocol.NewBulkString("first value"),
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
		"setting value with XX option does not set the value if key does not exist": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("missing-key-with-value-xx" + uniqueSuffix),
						protocol.NewBulkString("first value"),
						protocol.NewBulkString("XX"),
					},
					nil,
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("missing-key-with-value-xx" + uniqueSuffix),
					},
					nil,
				),
			},
		},
		"setting value with XX option sets the value if key has a current value": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-with-value-xx" + uniqueSuffix),
						protocol.NewBulkString("first value"),
					},
					protocol.NewSimpleString("OK"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-with-value-xx" + uniqueSuffix),
						protocol.NewBulkString("second value"),
						protocol.NewBulkString("XX"),
					},
					protocol.NewSimpleString("OK"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("key-with-value-xx" + uniqueSuffix),
					},
					protocol.NewBulkString("second value"),
				),
			},
		},
		"setting value with GEt and XX option return old value and sets the value if key has a current value": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-with-value-xx" + uniqueSuffix),
						protocol.NewBulkString("first value"),
					},
					protocol.NewSimpleString("OK"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-with-value-xx" + uniqueSuffix),
						protocol.NewBulkString("second value"),
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("XX"),
					},
					protocol.NewBulkString("first value"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("key-with-value-xx" + uniqueSuffix),
					},
					protocol.NewBulkString("second value"),
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
