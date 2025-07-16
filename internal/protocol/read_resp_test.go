package protocol_test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"redis-challenge/internal/protocol"
	"testing"
)

func TestParseBuffer(t *testing.T) {

	tests := map[string]struct {
		input         string
		expectedData  protocol.Data
		expectedBytes int
	}{
		"partial frame for a simple string": {
			input:         "+mess",
			expectedData:  nil,
			expectedBytes: 0,
		},
		"complete frame for a simple string": {
			input:         "+message\r\n",
			expectedData:  protocol.NewSimpleString("message"),
			expectedBytes: 7 + 3,
		},
		"complete frame for a simple string with partial of next frame": {
			input:         "+message\r\n+next",
			expectedData:  protocol.NewSimpleString("message"),
			expectedBytes: 7 + 3,
		},
		"partial frame for an error": {
			input:         "-error",
			expectedData:  nil,
			expectedBytes: 0,
		},
		"complete frame for an error": {
			input:         "-error\r\n",
			expectedData:  protocol.NewSimpleError("error"),
			expectedBytes: 5 + 3,
		},
		"complete frame for an error with partial of next frame": {
			input:         "-error\r\n+next",
			expectedData:  protocol.NewSimpleError("error"),
			expectedBytes: 5 + 3,
		},
		"partial frame for an integer": {
			input:         ":100",
			expectedData:  nil,
			expectedBytes: 0,
		},
		"complete frame for an integer": {
			input:         ":100\r\n",
			expectedData:  protocol.NewSimpleInteger(100),
			expectedBytes: 3 + 3,
		},
		"complete frame for a maximum positive integer": {
			input:         ":9223372036854775807\r\n",
			expectedData:  protocol.NewSimpleInteger(9223372036854775807),
			expectedBytes: 19 + 3,
		},
		"complete frame for a maximum negative integer": {
			input:         ":-9223372036854775806\r\n",
			expectedData:  protocol.NewSimpleInteger(-9223372036854775806),
			expectedBytes: 20 + 3,
		},
		"complete frame for an integer that is a floating-point number": {
			input:         ":1.25\r\n",
			expectedData:  protocol.NewSimpleError("value \"1.25\" is not a 64-bit integer"),
			expectedBytes: 4 + 3,
		},
		"complete frame for an integer that is too big": {
			input:         ":9223372036854775808\r\n",
			expectedData:  protocol.NewSimpleError("value \"9223372036854775808\" is not a 64-bit integer"),
			expectedBytes: len(":9223372036854775808\r\n"),
		},
		"complete frame for an integer with partial of next frame": {
			input:         ":100\r\n:99",
			expectedData:  protocol.NewSimpleInteger(100),
			expectedBytes: 3 + 3,
		},
		"complete frame for a null string": {
			input:         "$-1\r\n",
			expectedData:  nil,
			expectedBytes: 5,
		},
		"frame with partial size of a bulk string": {
			input:         "$10",
			expectedData:  nil,
			expectedBytes: 0,
		},
		"frame for a bulk string that is sized": {
			input:         "$10\r\n",
			expectedData:  nil,
			expectedBytes: 0,
		},
		"frame for a bulk string that is sized and is followed by partial of string": {
			input:         "$10\r\nabc",
			expectedData:  nil,
			expectedBytes: 0,
		},
		"frame for a bulk string that is sized and is followed by unterminated string of expected length": {
			input:         "$10\r\n1234567890",
			expectedData:  nil,
			expectedBytes: 0,
		},
		"frame for a bulk string that is sized and is followed by terminated string of expected length": {
			input:         "$10\r\n1234567890\r\n",
			expectedData:  protocol.NewBulkString("1234567890"),
			expectedBytes: len("$10\r\n1234567890\r\n"),
		},
		"frame for a bulk string that is sized with complete bulk string text and followed by a partial frame": {
			input:         "$3\r\nabc\r\n$10",
			expectedData:  protocol.NewBulkString("abc"),
			expectedBytes: len("$3\r\nabc\r\n"),
		},
		"frame for a bulk string cannot have a length that is a floating-point number": {
			input:         "$1.25\r\n",
			expectedData:  protocol.NewSimpleError("value \"1.25\" is not a valid bulk string length"),
			expectedBytes: 4 + 3,
		},
		"frame for an empty array": {
			input:         "*0\r\n",
			expectedData:  protocol.NewArray(nil),
			expectedBytes: 1 + 3,
		},
		"partial frame for an array": {
			input:         "*0",
			expectedData:  nil,
			expectedBytes: 0,
		},
		"complete frame for an array of size 1 with partial of next frame": {
			input:         "*1\r\n:10",
			expectedData:  nil,
			expectedBytes: 0,
		},
		"complete frame for an array of size 1 followed by completed simple integer frame": {
			input: "*1\r\n:10\r\n",
			expectedData: protocol.NewArray(
				[]protocol.Data{
					protocol.NewSimpleInteger(10),
				},
			),
			expectedBytes: len("*1\r\n:10\r\n"),
		},
		"complete frame for an array of size 2 followed by 2 complete frames with different types": {
			input: "*2\r\n:10\r\n+text\r\n",
			expectedData: protocol.NewArray(
				[]protocol.Data{
					protocol.NewSimpleInteger(10),
					protocol.NewSimpleString("text"),
				},
			),
			expectedBytes: len("*2\r\n:10\r\n+text\r\n"),
		},
		"complete frame for an array of size 2 followed by only 1 complete frame": {
			input:         "*2\r\n:10\r\n",
			expectedData:  nil,
			expectedBytes: 0,
		},
		"complete frame for an array of size 1 followed by an invalid completed simple integer frame": {
			input:         "*1\r\n:not-a-number\r\n",
			expectedData:  protocol.NewSimpleError("value \"not-a-number\" is not a 64-bit integer"),
			expectedBytes: len("*1\r\n:not-a-number\r\n"),
		},
		"frame with an unknown prefix": {
			input:         "xyz\r\n",
			expectedData:  protocol.NewSimpleError("unknown protocol symbol \"x\""),
			expectedBytes: 5,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var buffer bytes.Buffer
			buffer.WriteString(tt.input)

			data, byteCount := protocol.ReadFrame(&buffer)

			require.Equal(t, tt.expectedData, data)
			assert.Equal(t, tt.expectedBytes, byteCount)
		})
	}
}
