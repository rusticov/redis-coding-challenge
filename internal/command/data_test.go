package command_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"redis-challenge/internal/command"
	"redis-challenge/internal/protocol"
	"testing"
)

func TestReadingCommands(t *testing.T) {

	t.Run("simple ping command", func(t *testing.T) {
		data := protocol.Array{Data: []protocol.Data{
			protocol.NewBulkString("PING"),
		}}

		commandData, err := command.FromData(data)
		require.NoError(t, err)

		assert.Equal(t, "PING", commandData.Name)
		assert.Empty(t, commandData.Arguments)
	})

	t.Run("ping command with message", func(t *testing.T) {
		data := protocol.Array{Data: []protocol.Data{
			protocol.NewBulkString("PING"),
			protocol.NewBulkString("message"),
		}}

		commandData, err := command.FromData(data)
		require.NoError(t, err)

		assert.Equal(t, "PING", commandData.Name)
		assert.Equal(t, []protocol.Data{
			protocol.NewBulkString("message"),
		}, commandData.Arguments)
	})

	t.Run("echo command with message", func(t *testing.T) {
		data := protocol.Array{Data: []protocol.Data{
			protocol.NewBulkString("ECHO"),
			protocol.NewBulkString("echo message"),
		}}

		commandData, err := command.FromData(data)
		require.NoError(t, err)

		assert.Equal(t, "ECHO", commandData.Name)
		assert.Equal(t, []protocol.Data{
			protocol.NewBulkString("echo message"),
		}, commandData.Arguments)
	})

	t.Run("data is an empty array should error", func(t *testing.T) {
		data := protocol.Array{}

		_, err := command.FromData(data)
		require.EqualError(t, err, "missing command name")
	})

	t.Run("data whose name is not a bulk string should error", func(t *testing.T) {
		data := protocol.Array{Data: []protocol.Data{
			protocol.NewSimpleString("PING"),
		}}

		_, err := command.FromData(data)
		require.EqualError(t, err, "command name must be a bulk string")
	})

	t.Run("data is not an array should error", func(t *testing.T) {
		data := protocol.NewBulkString("PING")

		_, err := command.FromData(data)
		require.EqualError(t, err, "not a command")
	})
}
