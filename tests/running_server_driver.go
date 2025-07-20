package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"redis-challenge/internal/server"
	"redis-challenge/tests/call"
	"testing"
	"time"
)

const (
	timeout              = 100 * time.Millisecond
	LargeStringByteCount = 8192
)

type ServerVariant string

const (
	UseRealRedisServer ServerVariant = "real redis server"
	UseChallengeServer ServerVariant = ""
)

func DriveProtocolAgainstServer[T call.Call](t testing.TB, calls []T, variant ...ServerVariant) {
	testServer := createTestServer(t, variant...)
	defer func() {
		require.NoError(t, testServer.Close(), "failed to close test server")
	}()

	connection, err := net.DialTimeout("tcp", testServer.Address(), timeout)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, connection.Close(), "failed to close connection to the test server")
	}()

	for _, call := range calls {
		request := call.Request()
		_, err = connection.Write([]byte(request))
		require.NoError(t, err, "failed to write request: %s", request)

		if !call.IsResponseExpected() {
			continue
		}

		buffer := make([]byte, LargeStringByteCount+20)
		n, err := connection.Read(buffer)
		assert.NoError(t, err, "failed to read reply to the request: %s", request)

		response := string(buffer[:n])

		call.ConfirmResponse(t, response)
		//assert.Equal(t, call.expectedResponse, response,
		//	"unexpected reply to the request: %s", call.request)
	}
}

func createTestServer(t testing.TB, variant ...ServerVariant) server.Server {
	activeVariant := UseChallengeServer
	if len(variant) > 0 {
		activeVariant = variant[0]
	}

	switch activeVariant {
	case UseRealRedisServer:
		return NewRealRedisServer()
	default:
		challengeServer, err := server.NewChallengeServer()
		require.NoError(t, err)
		return challengeServer
	}
}
