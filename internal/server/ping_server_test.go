package server_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"redis-challenge/internal/server"
	"testing"
	"time"
)

func TestPingServer(t *testing.T) {

	timeout := 100 * time.Millisecond

	tests := map[string]struct {
		command          string
		expectedResponse string
	}{
		"send ping without message and receive PONG": {
			command:          "*1\r\n$4\r\nPING\r\n",
			expectedResponse: "+PONG\r\n",
		},
		"send ping with message should receive message back in reply": {
			command:          "*2\r\n$4\r\nPING\r\n$11\r\nthe message\r\n",
			expectedResponse: "$11\r\nthe message\r\n",
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			testServer := createTestServer(t)
			defer testServer.Close()

			connection, err := net.DialTimeout("tcp", testServer.Address(), timeout)
			require.NoError(t, err)
			defer connection.Close()

			_, err = connection.Write([]byte(test.command))
			require.NoError(t, err)

			buffer := make([]byte, 256)
			n, err := connection.Read(buffer)
			assert.NoError(t, err)

			response := string(buffer[:n])
			assert.Equal(t, test.expectedResponse, response)
		})
	}
}

type ServerVariant string

const (
	RealRedisServer ServerVariant = "real redis server"
	ChallengeServer ServerVariant = "challenge server"
)

func createTestServer(t testing.TB, variant ...ServerVariant) server.Server {
	activeVariant := ChallengeServer
	if len(variant) > 0 {
		activeVariant = variant[0]
	}

	switch activeVariant {
	case RealRedisServer:
		return server.NewRealRedisServer()
	default:
		challengeServer, err := server.NewChallengeServer()
		require.NoError(t, err)
		return challengeServer
	}
}
