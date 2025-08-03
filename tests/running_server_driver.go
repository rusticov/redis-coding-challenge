package tests

import (
	"github.com/stretchr/testify/require"
	"io"
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

func (s ServerVariant) Sleep(clock store.Clock, c call.Call) {
	delay := c.Delay()

	if s == UseChallengeServer {
		if fixedClock, ok := clock.(*store.FixedClock); ok {
			fixedClock.AddMilliseconds(delay.Milliseconds())
			return
		}
	}

	time.Sleep(delay)
}

func DriveProtocolAgainstServer[T call.Call](t testing.TB, calls []T, variant ServerVariant, options ...any) {
	var clock store.Clock = &store.FixedClock{TimeInMilliseconds: time.Now().UnixMilli()}
	logWriter := io.Discard

	for _, option := range options {
		if definedClock, ok := option.(io.Writer); ok {
			logWriter = definedClock
		}

		if definedClock, ok := option.(store.Clock); ok {
			clock = definedClock
		}
	}

	testServer := createTestServer(t, clock, variant, logWriter)
	defer func() {
		require.NoError(t, testServer.Close(), "failed to close test server")
	}()

	SendCallsToServer(t, testServer, calls, variant, clock)
}

func SendCallsToServer[T call.Call](t testing.TB, testServer server.Server, calls []T, variant ServerVariant, clock store.Clock) {

	connection, err := net.DialTimeout("tcp", testServer.Address(), timeout)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, connection.Close(), "failed to close connection to the test server")
	}()

	for _, nextCall := range calls {
		variant.Sleep(clock, nextCall)

		request := nextCall.Request()
		_, err := connection.Write([]byte(request))
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

func createTestServer(t testing.TB, clock store.Clock, variant ServerVariant, logWriter io.Writer) server.Server {
	switch variant {
	case UseRealRedisServer:
		return NewRealRedisServer()
	default:
		challengeServer, err := server.NewChallengeServer(0, store.NewBuilder()).
			WithClock(clock).
			WithArchiveWriter(logWriter).
			Start()
		require.NoError(t, err)
		return challengeServer
	}
}
