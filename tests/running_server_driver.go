package tests

import (
	"github.com/stretchr/testify/require"
	"net"
	"redis-challenge/internal/server"
	"redis-challenge/internal/store"
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

	for _, nextCall := range calls {
		time.Sleep(nextCall.Delay())

		request := nextCall.Request()
		_, err = connection.Write([]byte(request))
		require.NoError(t, err, "failed to write request: %s", request)

		if !nextCall.IsResponseExpected() {
			continue
		}

		response := ""

		buffer := make([]byte, LargeStringByteCount+20)

		for nextCall.IsPossiblePartialResponse(response) {
			n, err := connection.Read(buffer)
			require.NoError(t, err, "failed to read reply to the request: %s", request)

			response += string(buffer[:n])
		}

		nextCall.ConfirmResponse(t, response)
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
		challengeServer, err := server.NewChallengeServer(0, store.New())
		require.NoError(t, err)
		return challengeServer
	}
}
