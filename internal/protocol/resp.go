package protocol

import (
	"bytes"
	"fmt"
	"strconv"
)

func ReadFrame(b *bytes.Buffer) (Data, int) {
	bs := b.Bytes()
	delimiterIndex := bytes.Index(bs, []byte("\r\n"))
	if delimiterIndex == -1 {
		return nil, 0
	}

	text := string(bs[1:delimiterIndex])

	symbol := bs[0]

	frameSize := delimiterIndex + 2

	switch symbol {
	case '-':
		return NewSimpleError(text), frameSize
	case ':':
		return parseSimpleInteger(text, frameSize)
	case '$':
		return parseBulkString(text, frameSize, bs)
	case '*':
		return NewArray(0), frameSize
	case '+':
		return NewSimpleString(text), frameSize
	default:
		return NewSimpleError(fmt.Sprintf("unknown protocol symbol \"%c\"", symbol)), frameSize
	}
}

func parseSimpleInteger(text string, frameSize int) (Data, int) {
	value, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return NewSimpleError(fmt.Sprintf("value \"%s\" is not an integer", text)), frameSize
	}
	return NewSimpleInteger(value), frameSize
}

func parseBulkString(text string, frameSize int, bs []byte) (Data, int) {
	if text == "-1" {
		return nil, 5
	}

	length, err := strconv.Atoi(text)
	if err != nil {
		return NewSimpleError(fmt.Sprintf("value \"%s\" is not a valid bulk string length", text)), frameSize
	}

	if frameSize+length+2 <= len(bs) {
		return NewBulkString(string(bs[frameSize : frameSize+length])), frameSize + length + 2
	}
	return nil, 0
}
