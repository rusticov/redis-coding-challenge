package protocol

import (
	"bytes"
	"fmt"
	"strconv"
)

func ReadFrame(bs []byte) (Data, int) {
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
		return parseArray(bs, text, frameSize)
	case '+':
		return NewSimpleString(text), frameSize
	default:
		return NewSimpleError(fmt.Sprintf("unknown protocol symbol \"%c\"", symbol)), frameSize
	}
}

func parseSimpleInteger(text string, frameSize int) (Data, int) {
	value, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return NewSimpleError(fmt.Sprintf("value \"%s\" is not a 64-bit integer", text)), frameSize
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

func parseArray(bs []byte, text string, frameSize int) (Data, int) {
	length, _ := strconv.Atoi(text)

	if length == 0 {
		return NewArray(nil), frameSize
	}

	data := make([]Data, length)
	for i := range length {
		datum, datumSize := ReadFrame(bs[frameSize:])

		if datum == nil {
			return nil, 0
		}

		data[i] = datum
		frameSize += datumSize

		if errorData, ok := datum.(SimpleError); ok {
			return errorData, frameSize
		}
	}
	return NewArray(data), frameSize
}
