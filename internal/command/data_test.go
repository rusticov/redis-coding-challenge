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
		require.Nil(t, err, "should be no error")

		assert.Equal(t, "PING", commandData.Name)
		assert.Empty(t, commandData.Arguments)
	})

	t.Run("ping command with message", func(t *testing.T) {
		data := protocol.Array{Data: []protocol.Data{
			protocol.NewBulkString("PING"),
			protocol.NewBulkString("message"),
		}}

		commandData, err := command.FromData(data)
		require.Nil(t, err, "should be no error")

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
		require.Nil(t, err, "should be no error")

		assert.Equal(t, "ECHO", commandData.Name)
		assert.Equal(t, []protocol.Data{
			protocol.NewBulkString("echo message"),
		}, commandData.Arguments)
	})

	t.Run("data is an empty array should error", func(t *testing.T) {
		data := protocol.Array{}

		_, err := command.FromData(data)
		assert.Equal(t, protocol.NewSimpleError("missing command name"), err)
	})

	t.Run("data whose name is a simple string should error", func(t *testing.T) {
		data := protocol.Array{Data: []protocol.Data{
			protocol.NewSimpleString("PING"),
		}}

		_, err := command.FromData(data)
		assert.Equal(t, protocol.NewSimpleError("ERR Protocol error: expected '$', got '+'"), err)
	})

	t.Run("data whose name is a non-string should error", func(t *testing.T) {
		data := protocol.Array{Data: []protocol.Data{
			protocol.NewSimpleInteger(1),
		}}

		_, err := command.FromData(data)
		assert.Equal(t, protocol.NewSimpleError("ERR Protocol error: expected '$', got ':'"), err)
	})

	t.Run("data is not an array should error", func(t *testing.T) {
		data := protocol.NewBulkString("PING")

		_, err := command.FromData(data)
		assert.Equal(t, protocol.NewSimpleError("not a command"), err)
	})

	t.Run("data is an error should be returned as the error", func(t *testing.T) {
		data := protocol.NewSimpleError("I am an error")

		_, err := command.FromData(data)
		assert.Equal(t, protocol.NewSimpleError("I am an error"), err)
	})
}
