package server_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"redis-challenge/internal/server"
	"testing"
	"time"
)

type CallToRedis struct {
	request          string
	expectedResponse string
}

func TestPingServer(t *testing.T) {

	timeout := 100 * time.Millisecond

	tests := map[string]struct {
		calls   []CallToRedis
		variant ServerVariant
	}{
		"send ping without message and receive PONG": {
			calls: []CallToRedis{{
				request:          "*1\r\n$4\r\nPING\r\n",
				expectedResponse: "+PONG\r\n",
			}},
		},
		"send ping with message should receive message back in reply": {
			calls: []CallToRedis{{
				request:          "*2\r\n$4\r\nPING\r\n$11\r\nthe message\r\n",
				expectedResponse: "$11\r\nthe message\r\n",
			}},
		},
		"send echo with message should receive message back in reply": {
			calls: []CallToRedis{{
				request:          "*2\r\n$4\r\nECHO\r\n$12\r\necho message\r\n",
				expectedResponse: "$12\r\necho message\r\n",
			}},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			testServer := createTestServer(t, test.variant)
			defer func() {
				err := testServer.Close()
				require.NoError(t, err, "failed to close test server")
			}()

			connection, err := net.DialTimeout("tcp", testServer.Address(), timeout)
			require.NoError(t, err)
			defer func() {
				err := connection.Close()
				require.NoError(t, err, "failed to close connection to the test server")
			}()

			for _, call := range test.calls {
				_, err = connection.Write([]byte(call.request))
				require.NoError(t, err, "failed to write request: %s", call.request)

				buffer := make([]byte, 256)
				n, err := connection.Read(buffer)
				assert.NoError(t, err, "failed to read reply to the request: %s", call.request)

				response := string(buffer[:n])
				assert.Equal(t, call.expectedResponse, response,
					"unexpected reply to the request: %s", call.request)
			}
		})
	}

	t.Run("send bad command should receive error message", func(t *testing.T) {
		testServer := createTestServer(t)
		defer func() {
			err := testServer.Close()
			require.NoError(t, err, "failed to close test server")
		}()

		connection, err := net.DialTimeout("tcp", testServer.Address(), timeout)
		require.NoError(t, err)
		defer func() {
			err := connection.Close()
			require.NoError(t, err, "failed to close connection to the test server")
		}()

		_, err = connection.Write([]byte("*2\r\n$3\r\nBAD\r\n$3\r\narg\r\n"))
		require.NoError(t, err)

		buffer := make([]byte, 256)
		n, err := connection.Read(buffer)
		assert.NoError(t, err)

		response := string(buffer[:n])
		assert.Contains(t, "-ERR unknown command 'BAD'\r\n", response)
	})
}

type ServerVariant string

const (
	RealRedisServer ServerVariant = "real redis server"
	ChallengeServer ServerVariant = ""
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
