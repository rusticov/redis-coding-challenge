package protocol_test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"redis-challenge/internal/protocol"
	"testing"
)

func TestParseBuffer(t *testing.T) {

	tests := []struct {
		name          string
		input         string
		expectedData  protocol.Data
		expectedBytes int
	}{
		{
			name:          "partial frame for a simple string",
			input:         "+mess",
			expectedData:  nil,
			expectedBytes: 0,
		},
		{
			name:          "complete frame for a simple string",
			input:         "+message\r\n",
			expectedData:  protocol.NewSimpleString("message"),
			expectedBytes: 7 + 3,
		},
		{
			name:          "complete frame for a simple string with partial of next frame",
			input:         "+message\r\n+next",
			expectedData:  protocol.NewSimpleString("message"),
			expectedBytes: 7 + 3,
		},
		{
			name:          "partial frame for an error",
			input:         "-error",
			expectedData:  nil,
			expectedBytes: 0,
		},
		{
			name:          "complete frame for an error",
			input:         "-error\r\n",
			expectedData:  protocol.NewError("error"),
			expectedBytes: 5 + 3,
		},
		{
			name:          "complete frame for an error with partial of next frame",
			input:         "-error\r\n+next",
			expectedData:  protocol.NewError("error"),
			expectedBytes: 5 + 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buffer bytes.Buffer
			buffer.WriteString(tt.input)

			data, byteCount := protocol.ReadFrame(&buffer)

			assert.Equal(t, tt.expectedData, data)
			assert.Equal(t, tt.expectedBytes, byteCount)
		})
	}
}
