package protocol_test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"redis-challenge/internal/protocol"
	"testing"
)

func TestWritingData(t *testing.T) {

	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "simple string",
			message: "+message\r\n",
		},
		{
			name:    "simple error",
			message: "-simple error message\r\n",
		},
		{
			name:    "simple integer",
			message: ":42\r\n",
		},
		{
			name:    "bulk string",
			message: "$5\r\nabcde\r\n",
		},
		{
			name:    "array",
			message: "*2\r\n+abcde\r\n:42\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buffer bytes.Buffer
			buffer.WriteString(tt.message)

			data, _ := protocol.ReadFrame(&buffer)

			var outBuffer bytes.Buffer
			err := protocol.WriteData(&outBuffer, data)
			require.NoError(t, err)

			assert.Equal(t, tt.message, outBuffer.String())
		})
	}
}
