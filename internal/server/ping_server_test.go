package server_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"redis-challenge/internal/server"
	"strings"
	"testing"
	"time"
)

type CallToRedis struct {
	request          string
	expectedResponse string
}

func TestPingServer(t *testing.T) {

	timeout := 100 * time.Millisecond

	const largeStringByteCount = 8192

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
		"send two echos and receive replies to each": {
			calls: []CallToRedis{
				{
					request:          "*2\r\n$4\r\nECHO\r\n$5\r\nfirst\r\n",
					expectedResponse: "$5\r\nfirst\r\n",
				},
				{
					request:          "*2\r\n$4\r\nECHO\r\n$6\r\nsecond\r\n",
					expectedResponse: "$6\r\nsecond\r\n",
				},
			},
		},
		"send echo split across 2 requests": {
			calls: []CallToRedis{
				{
					request: "*2\r\n$4\r\nECHO\r\n$5\r\nfi",
				},
				{
					request:          "rst\r\n",
					expectedResponse: "$5\r\nfirst\r\n",
				},
			},
		},
		"send echo with large message": {
			calls: []CallToRedis{
				{
					request:          fmt.Sprintf("*2\r\n$4\r\nECHO\r\n$%d\r\n%s\r\n", largeStringByteCount, strings.Repeat("x", largeStringByteCount)),
					expectedResponse: fmt.Sprintf("$%d\r\n%s\r\n", largeStringByteCount, strings.Repeat("x", largeStringByteCount)),
				},
			},
		},
		"send echo no message should reply with error": {
			calls: []CallToRedis{{
				request:          "*1\r\n$4\r\nECHO\r\n",
				expectedResponse: "-ERR wrong number of arguments for 'echo' command\r\n",
			}},
		},
		"send empty array should reply with error": {
			calls: []CallToRedis{{
				request:          "*1\r\n+ECHO\r\n",
				expectedResponse: "-ERR Protocol error: expected '$', got '+'\r\n",
			}},
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			testServer := createTestServer(t, test.variant)
			defer func() {
				require.NoError(t, testServer.Close(), "failed to close test server")
			}()

			connection, err := net.DialTimeout("tcp", testServer.Address(), timeout)
			require.NoError(t, err)
			defer func() {
				require.NoError(t, connection.Close(), "failed to close connection to the test server")
			}()

			for _, call := range test.calls {
				_, err = connection.Write([]byte(call.request))
				require.NoError(t, err, "failed to write request: %s", call.request)

				if call.expectedResponse == "" {
					continue
				}

				buffer := make([]byte, largeStringByteCount+20)
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
			require.NoError(t, testServer.Close(), "failed to close test server")
		}()

		connection, err := net.DialTimeout("tcp", testServer.Address(), timeout)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, connection.Close(), "failed to close connection to the test server")
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
