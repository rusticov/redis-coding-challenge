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

	t.Run("send ping without message and receive PONG", func(t *testing.T) {
		testServer := createTestServer(t)
		defer testServer.Close()

		connection, err := net.DialTimeout("tcp", testServer.Address(), timeout)
		require.NoError(t, err)
		defer connection.Close()

		_, err = connection.Write([]byte("*1\r\n$4\r\nPING\r\n"))
		require.NoError(t, err)

		buffer := make([]byte, 256)
		n, err := connection.Read(buffer)
		assert.NoError(t, err)

		response := string(buffer[:n])
		assert.Equal(t, "+PONG\r\n", response)
	})

	t.Run("send ping with message should receive message back in reply", func(t *testing.T) {
		testServer := createTestServer(t, RealRedisServer)
		defer testServer.Close()

		connection, err := net.DialTimeout("tcp", testServer.Address(), timeout)
		require.NoError(t, err)
		defer connection.Close()

		_, err = connection.Write([]byte("*2\r\n$4\r\nPING\r\n$11\r\nthe message\r\n"))
		require.NoError(t, err)

		buffer := make([]byte, 256)
		n, err := connection.Read(buffer)
		assert.NoError(t, err)

		response := string(buffer[:n])
		assert.Equal(t, "$11\r\nthe message\r\n", response)
	})
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
