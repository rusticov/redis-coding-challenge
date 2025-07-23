package tests_test

import (
	nanoid "github.com/matoous/go-nanoid/v2"
	"redis-challenge/internal/protocol"
	"redis-challenge/tests"
	"redis-challenge/tests/call"
	"testing"
)

func TestDeletingValues(t *testing.T) {

	uniqueSuffix := "-" + nanoid.Must(6)

	testCases := executionTestCases{
		"deleting 1 value that has not been set should declare that 0 have been deleted": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("DEL"),
						protocol.NewBulkString("random-key-with-no-value" + uniqueSuffix),
					},
					protocol.NewSimpleInteger(0),
				),
			},
		},
		"deleting 1 value that has been set should declare 1 is deleted": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-to-delete" + uniqueSuffix),
						protocol.NewBulkString("value 1"),
					},
					protocol.NewSimpleString("OK"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("DEL"),
						protocol.NewBulkString("key-to-delete" + uniqueSuffix),
					},
					protocol.NewSimpleInteger(1),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("GET"),
						protocol.NewBulkString("key-to-delete" + uniqueSuffix),
					},
					nil,
				),
			},
		},
		"deleting values with mix is known and unknown key should return count of deleted keys that were known": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-to-delete-1" + uniqueSuffix),
						protocol.NewBulkString("value 1"),
					},
					protocol.NewSimpleString("OK"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-to-delete-3" + uniqueSuffix),
						protocol.NewBulkString("value 3"),
					},
					protocol.NewSimpleString("OK"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("DEL"),
						protocol.NewBulkString("key-to-delete-1" + uniqueSuffix),
						protocol.NewBulkString("key-to-delete-2" + uniqueSuffix),
						protocol.NewBulkString("key-to-delete-3" + uniqueSuffix),
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
