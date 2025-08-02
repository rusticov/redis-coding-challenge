package tests

import (
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"redis-challenge/internal/command"
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

func DriveProtocolAgainstServer[T call.Call](t testing.TB, calls []T, variant ServerVariant, logWriter ...io.Writer) {
	clock := &store.FixedClock{TimeInMilliseconds: time.Now().UnixMilli()}

	actualLogWriter := io.Discard
	if len(logWriter) > 0 {
		actualLogWriter = logWriter[0]
	}

	testServer := createTestServer(t, clock, variant, actualLogWriter)
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
		tracker := store.NewExpiryTracker()
		s := store.NewWithClock(clock).WithExpiryTracker(tracker)
		scanner := command.NewExpiryScanner(tracker, s)

		challengeServer, err := server.NewChallengeServer(0, s, scanner).WithWriter(logWriter).Start()
		require.NoError(t, err)
		return challengeServer
	}
}
