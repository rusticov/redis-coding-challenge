package command_test

import (
	nanoid "github.com/matoous/go-nanoid/v2"
	"redis-challenge/internal/protocol"
	"redis-challenge/tests"
	"redis-challenge/tests/call"
	"testing"
)

func TestMixSettingValuesOfIncompatibleTypes(t *testing.T) {

	uniqueSuffix := "-" + nanoid.Must(6)

	testCases := executionTestCases{
		"lpush to key with a string value should return error": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-set-lpush" + uniqueSuffix),
						protocol.NewBulkString("value 1"),
					},
					protocol.NewSimpleString("OK"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LPUSH"),
						protocol.NewBulkString("key-set-lpush" + uniqueSuffix),
						protocol.NewBulkString("value 1"),
					},
					protocol.NewSimpleError("WRONGTYPE Operation against a key holding the wrong kind of value"),
				),
			},
		},
		"rpush to key with a string value should return error": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-set-rpush" + uniqueSuffix),
						protocol.NewBulkString("value 1"),
					},
					protocol.NewSimpleString("OK"),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("RPUSH"),
						protocol.NewBulkString("key-set-rpush" + uniqueSuffix),
						protocol.NewBulkString("value 1"),
					},
					protocol.NewSimpleError("WRONGTYPE Operation against a key holding the wrong kind of value"),
				),
			},
		},
		"set to key with a list value is ok": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LPUSH"),
						protocol.NewBulkString("key-lpush-set" + uniqueSuffix),
						protocol.NewBulkString("value 1"),
						protocol.NewBulkString("value 2"),
					},
					protocol.NewSimpleInteger(2),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("SET"),
						protocol.NewBulkString("key-lpush-set" + uniqueSuffix),
						protocol.NewBulkString("value 1"),
					},
					protocol.SimpleString("OK"),
				),
			},
		},
		"incr value to key with a list value should return error": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LPUSH"),
						protocol.NewBulkString("key-lpush-incr" + uniqueSuffix),
						protocol.NewBulkString("value 1"),
					},
					protocol.NewSimpleInteger(1),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("INCR"),
						protocol.NewBulkString("key-lpush-incr" + uniqueSuffix),
					},
					protocol.NewSimpleError("WRONGTYPE Operation against a key holding the wrong kind of value"),
				),
			},
		},
		"decr value to key with a list value should return error": {
			calls: []call.DataCall{
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("LPUSH"),
						protocol.NewBulkString("key-lpush-decr" + uniqueSuffix),
						protocol.NewBulkString("value 1"),
					},
					protocol.NewSimpleInteger(1),
				),
				call.NewFromData(
					[]protocol.Data{
						protocol.NewBulkString("DECR"),
						protocol.NewBulkString("key-lpush-decr" + uniqueSuffix),
					},
					protocol.NewSimpleError("WRONGTYPE Operation against a key holding the wrong kind of value"),
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
