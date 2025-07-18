package protocol_test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"redis-challenge/internal/protocol"
	"testing"
)

func TestWritingData(t *testing.T) {

	tests := map[string]string{
		"simple string":  "+message\r\n",
		"simple error":   "-simple error message\r\n",
		"simple integer": ":42\r\n",
		"bulk string":    "$5\r\nabcde\r\n",
		"array":          "*2\r\n+abcde\r\n:42\r\n",
	}

	for testName, message := range tests {
		t.Run(testName, func(t *testing.T) {
			var buffer bytes.Buffer
			buffer.WriteString(message)

			data, _ := protocol.ReadFrame(buffer.Bytes())

			var outBuffer bytes.Buffer
			err := protocol.WriteData(&outBuffer, data)
			require.NoError(t, err)

			assert.Equal(t, message, outBuffer.String())
		})
	}
}
