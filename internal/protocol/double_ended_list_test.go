package protocol_test

import (
	"bytes"
	"redis-challenge/internal/list"
	"redis-challenge/internal/protocol"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodingDoubleEndedList(t *testing.T) {

	t.Run("encoding a DoubleEndedList can be decoded as a list of bulk strings", func(t *testing.T) {
		doubleEndedList, ok := list.LeftPush([]string{"a", "b"}, nil)
		require.True(t, ok)
		doubleEndedList, ok = list.RightPush([]string{"c", "d"}, doubleEndedList)
		require.True(t, ok)

		buffer := bytes.NewBuffer(nil)
		err := protocol.WriteData(buffer, protocol.NewDoubleEndedList(doubleEndedList))
		require.NoError(t, err)

		decodedData, _ := protocol.ReadFrame(buffer.Bytes())
		assert.Equal(t, protocol.Array{Data: []protocol.Data{
			protocol.NewBulkString("b"),
			protocol.NewBulkString("a"),
			protocol.NewBulkString("c"),
			protocol.NewBulkString("d"),
		}}, decodedData)
	})
}
