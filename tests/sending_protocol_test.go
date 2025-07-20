package tests_test

import (
	"fmt"
	"redis-challenge/tests"
	"redis-challenge/tests/call"
	"strings"
	"testing"
)

func TestSendingProtocolToServer(t *testing.T) {

	const largeStringByteCount = tests.LargeStringByteCount

	testCases := map[string]struct {
		calls   []call.Call
		variant tests.ServerVariant
	}{
		"send ping without message and receive PONG": {
			calls: []call.Call{
				call.NewFromProtocol(
					"*1\r\n$4\r\nPING\r\n",
					"+PONG\r\n",
				),
			},
		},
		"send ping with message should receive message back in reply": {
			calls: []call.Call{
				call.NewFromProtocol(
					"*2\r\n$4\r\nPING\r\n$11\r\nthe message\r\n",
					"$11\r\nthe message\r\n",
				),
			},
		},
		"send echo with message should receive message back in reply": {
			calls: []call.Call{
				call.NewFromProtocol(
					"*2\r\n$4\r\nECHO\r\n$12\r\necho message\r\n",
					"$12\r\necho message\r\n",
				),
			},
		},
		"send two echos and receive replies to each": {
			calls: []call.Call{
				call.NewFromProtocol(
					"*2\r\n$4\r\nECHO\r\n$5\r\nfirst\r\n",
					"$5\r\nfirst\r\n",
				),
				call.NewFromProtocol(
					"*2\r\n$4\r\nECHO\r\n$6\r\nsecond\r\n",
					"$6\r\nsecond\r\n",
				),
			},
		},
		"send echo split across 2 requests": {
			calls: []call.Call{
				call.NewFromProtocolWithoutResponse(
					"*2\r\n$4\r\nECHO\r\n$5\r\nfi",
				),
				call.NewFromProtocol(
					"rst\r\n",
					"$5\r\nfirst\r\n",
				),
			},
		},
		"send echo with large message": {
			calls: []call.Call{
				call.NewFromProtocol(
					fmt.Sprintf("*2\r\n$4\r\nECHO\r\n$%d\r\n%s\r\n", largeStringByteCount, strings.Repeat("x", largeStringByteCount)),
					fmt.Sprintf("$%d\r\n%s\r\n", largeStringByteCount, strings.Repeat("x", largeStringByteCount)),
				),
			},
		},
		"send echo no message should reply with error": {
			calls: []call.Call{
				call.NewFromProtocol(
					"*1\r\n$4\r\nECHO\r\n",
					"-ERR wrong number of arguments for 'echo' command\r\n",
				),
			},
		},
		"send echo with extra argument should reply with error": {
			calls: []call.Call{
				call.NewFromProtocol(
					"*3\r\n$4\r\nECHO\r\n$3\r\none\r\n$3\r\ntwo\r\n",
					"-ERR wrong number of arguments for 'echo' command\r\n",
				),
			},
		},
		"send empty array should reply with error": {
			calls: []call.Call{
				call.NewFromProtocol(
					"*1\r\n+ECHO\r\n",
					"-ERR Protocol error: expected '$', got '+'\r\n",
				),
			},
		},
		"send bad command should receive error message": {
			calls: []call.Call{
				call.NewFromProtocolWithPartialResponse(
					"*2\r\n$3\r\nBAD\r\n$3\r\narg\r\n",
					"-ERR unknown command 'BAD'",
				),
			},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			tests.DriveProtocolAgainstServer(t, testCase.calls, testCase.variant)
		})
	}
}
