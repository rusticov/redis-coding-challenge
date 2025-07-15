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
			expectedData:  protocol.NewSimpleError("error"),
			expectedBytes: 5 + 3,
		},
		{
			name:          "complete frame for an error with partial of next frame",
			input:         "-error\r\n+next",
			expectedData:  protocol.NewSimpleError("error"),
			expectedBytes: 5 + 3,
		},
		{
			name:          "partial frame for an integer",
			input:         ":100",
			expectedData:  nil,
			expectedBytes: 0,
		},
		{
			name:          "complete frame for an integer",
			input:         ":100\r\n",
			expectedData:  protocol.NewSimpleInteger(100),
			expectedBytes: 3 + 3,
		},
		{
			name:          "complete frame for a maximum positive integer",
			input:         ":9223372036854775807\r\n",
			expectedData:  protocol.NewSimpleInteger(9223372036854775807),
			expectedBytes: 19 + 3,
		},
		{
			name:          "complete frame for a maximum negative integer",
			input:         ":-9223372036854775806\r\n",
			expectedData:  protocol.NewSimpleInteger(-9223372036854775806),
			expectedBytes: 20 + 3,
		},
		{
			name:          "complete frame for an integer that is a floating-point number",
			input:         ":1.25\r\n",
			expectedData:  protocol.NewSimpleError("value \"1.25\" is not an integer"),
			expectedBytes: 4 + 3,
		},
		{
			name:          "complete frame for an integer that is too big",
			input:         ":9223372036854775808\r\n",
			expectedData:  protocol.NewSimpleError("value \"9223372036854775808\" is not an integer"),
			expectedBytes: len(":9223372036854775808\r\n"),
		},
		{
			name:          "complete frame for an integer with partial of next frame",
			input:         ":100\r\n:99",
			expectedData:  protocol.NewSimpleInteger(100),
			expectedBytes: 3 + 3,
		},
		{
			name:          "complete frame for a null string",
			input:         "$-1\r\n",
			expectedData:  nil,
			expectedBytes: 5,
		},
		{
			name:          "frame with partial size of a bulk string",
			input:         "$10",
			expectedData:  nil,
			expectedBytes: 0,
		},
		{
			name:          "frame for a bulk string that is sized",
			input:         "$10\r\n",
			expectedData:  protocol.NewBulkStringStart(10),
			expectedBytes: len("$10\r\n"),
		},
		{
			name:          "frame for a bulk string that is sized and is followed by partial of next frame",
			input:         "$10\r\nabc",
			expectedData:  protocol.NewBulkStringStart(10),
			expectedBytes: len("$10\r\n"),
		},
		{
			name:          "frame for a bulk string cannot have a length that is a floating-point number",
			input:         "$1.25\r\n",
			expectedData:  protocol.NewSimpleError("value \"1.25\" is not a valid bulk string length"),
			expectedBytes: 4 + 3,
		},
		{
			name:          "frame for an empty array",
			input:         "*0\r\n",
			expectedData:  protocol.NewArray(0),
			expectedBytes: 1 + 3,
		},
		{
			name:          "partial frame for an array",
			input:         "*0",
			expectedData:  nil,
			expectedBytes: 0,
		},
		{
			name:          "frame with an unknown prefix",
			input:         "xyz\r\n",
			expectedData:  protocol.NewSimpleError("unknown protocol symbol \"x\""),
			expectedBytes: 5,
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
