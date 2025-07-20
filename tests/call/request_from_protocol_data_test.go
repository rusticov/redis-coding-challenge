package call_test

import (
	"github.com/stretchr/testify/assert"
	"redis-challenge/internal/protocol"
	"redis-challenge/tests/call"
	"testing"
)

func TestRequestFromProtocolData(t *testing.T) {

	t.Run("parse simple ping command into resp protocol request string", func(t *testing.T) {
		dataCall := call.NewFromDataWithoutError(
			[]protocol.Data{
				protocol.NewBulkString("PING"),
			},
		)

		request := dataCall.Request()

		assert.Equal(t, "*1\r\n$4\r\nPING\r\n", request)
	})

	t.Run("request data is wrapped in an array", func(t *testing.T) {
		dataCall := call.NewFromDataWithoutError(
			[]protocol.Data{
				protocol.NewBulkString("PING"),
			},
		)

		requestData := dataCall.RequestData()

		assert.Equal(t, protocol.NewArray(
			[]protocol.Data{
				protocol.NewBulkString("PING"),
			},
		), requestData)
	})
}
