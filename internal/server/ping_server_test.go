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

	timeout := 10 * time.Millisecond

	t.Run("send ping without message and receive PONG", func(t *testing.T) {
		testServer := createTestServer(t)

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
}

func createTestServer(t testing.TB) server.Server {
	testServer, err := server.NewRealRedisServer()
	require.NoError(t, err)

	return testServer
}
