package command_test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"redis-challenge/internal/command"
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
	"testing"
)

func TestNormalisingCommandSetWithExpiry(t *testing.T) {

	t.Run("set command with relative expiry EX normalises to PXAT (absolute expiry in milliseconds)", func(t *testing.T) {
		// Given a validator at a known time
		clock := &store.FixedClock{TimeInMilliseconds: 10_123}
		validator := command.NewValidator(clock)

		// Given a request with only the EX option
		request := protocol.Array{Data: []protocol.Data{
			protocol.NewBulkString("SET"),
			protocol.NewBulkString("key"),
			protocol.NewBulkString("value"),
			protocol.NewBulkString("EX"),
			protocol.NewBulkString("3"),
		}}

		buffer := bytes.NewBuffer(nil)
		err := protocol.WriteData(buffer, request)
		require.NoError(t, err)

		// When we validate a request with the EX option
		cmd, errorData := validator.Validate(buffer.Bytes(), request)

		require.Nil(t, errorData)
		requestBytes, requestType := cmd.Request()

		// Then the command is of an updating type
		assert.Equal(t, command.TypeUpdate, requestType)

		// Then the expiry is converted to absolute time in milliseconds
		validatedRequest, _ := protocol.ReadFrame(requestBytes)
		assert.Equal(t, protocol.Array{Data: []protocol.Data{
			protocol.NewBulkString("SET"),
			protocol.NewBulkString("key"),
			protocol.NewBulkString("value"),
			protocol.NewBulkString("PXAT"),
			protocol.NewBulkString("13123"),
		}}, validatedRequest)
	})

	t.Run("set command with relative expiry PX normalises to PXAT (absolute expiry in milliseconds)", func(t *testing.T) {
		// Given a validator at a known time
		clock := &store.FixedClock{TimeInMilliseconds: 10_123}
		validator := command.NewValidator(clock)

		// Given a request with only the EX option
		request := protocol.Array{Data: []protocol.Data{
			protocol.NewBulkString("SET"),
			protocol.NewBulkString("key"),
			protocol.NewBulkString("value"),
			protocol.NewBulkString("PX"),
			protocol.NewBulkString("3456"),
		}}

		buffer := bytes.NewBuffer(nil)
		err := protocol.WriteData(buffer, request)
		require.NoError(t, err)

		// When we validate a request with the EX option
		cmd, errorData := validator.Validate(buffer.Bytes(), request)

		require.Nil(t, errorData)
		requestBytes, requestType := cmd.Request()

		// Then the command is of an updating type
		assert.Equal(t, command.TypeUpdate, requestType)

		// Then the expiry is converted to absolute time in milliseconds
		validatedRequest, _ := protocol.ReadFrame(requestBytes)
		assert.Equal(t, protocol.Array{Data: []protocol.Data{
			protocol.NewBulkString("SET"),
			protocol.NewBulkString("key"),
			protocol.NewBulkString("value"),
			protocol.NewBulkString("PXAT"),
			protocol.NewBulkString("13579"),
		}}, validatedRequest)
	})

	t.Run("set command with relative expiry EX normalises to PXAT with preceding XX option", func(t *testing.T) {
		// Given a validator at a known time
		clock := &store.FixedClock{TimeInMilliseconds: 10_123}
		validator := command.NewValidator(clock)

		// Given a request with only the EX option
		request := protocol.Array{Data: []protocol.Data{
			protocol.NewBulkString("SET"),
			protocol.NewBulkString("key"),
			protocol.NewBulkString("value"),
			protocol.NewBulkString("XX"),
			protocol.NewBulkString("EX"),
			protocol.NewBulkString("3"),
		}}

		buffer := bytes.NewBuffer(nil)
		err := protocol.WriteData(buffer, request)
		require.NoError(t, err)

		// When we validate a request with the EX option
		cmd, errorData := validator.Validate(buffer.Bytes(), request)

		require.Nil(t, errorData)
		requestBytes, requestType := cmd.Request()

		// Then the command is of an updating type
		assert.Equal(t, command.TypeUpdate, requestType)

		// Then the expiry is converted to absolute time in milliseconds
		validatedRequest, _ := protocol.ReadFrame(requestBytes)
		assert.Equal(t, protocol.Array{Data: []protocol.Data{
			protocol.NewBulkString("SET"),
			protocol.NewBulkString("key"),
			protocol.NewBulkString("value"),
			protocol.NewBulkString("XX"),
			protocol.NewBulkString("PXAT"),
			protocol.NewBulkString("13123"),
		}}, validatedRequest)
	})

	t.Run("set command with relative expiry EX normalises to PXAT with following NX option", func(t *testing.T) {
		// Given a validator at a known time
		clock := &store.FixedClock{TimeInMilliseconds: 10_123}
		validator := command.NewValidator(clock)

		// Given a request with only the EX option
		request := protocol.Array{Data: []protocol.Data{
			protocol.NewBulkString("SET"),
			protocol.NewBulkString("key"),
			protocol.NewBulkString("value"),
			protocol.NewBulkString("EX"),
			protocol.NewBulkString("3"),
			protocol.NewBulkString("NX"),
		}}

		buffer := bytes.NewBuffer(nil)
		err := protocol.WriteData(buffer, request)
		require.NoError(t, err)

		// When we validate a request with the EX option
		cmd, errorData := validator.Validate(buffer.Bytes(), request)

		require.Nil(t, errorData)
		requestBytes, requestType := cmd.Request()

		// Then the command is of an updating type
		assert.Equal(t, command.TypeUpdate, requestType)

		// Then the expiry is converted to absolute time in milliseconds
		validatedRequest, _ := protocol.ReadFrame(requestBytes)
		assert.Equal(t, protocol.Array{Data: []protocol.Data{
			protocol.NewBulkString("SET"),
			protocol.NewBulkString("key"),
			protocol.NewBulkString("value"),
			protocol.NewBulkString("PXAT"),
			protocol.NewBulkString("13123"),
			protocol.NewBulkString("NX"),
		}}, validatedRequest)
	})
}
