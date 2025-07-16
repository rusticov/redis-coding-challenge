package protocol

import (
	"bytes"
	"fmt"
	"strconv"
)

type Data interface {
	IsData()
}

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
		value, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return NewSimpleError(fmt.Sprintf("value \"%s\" is not an integer", text)), frameSize
		}
		return NewSimpleInteger(value), frameSize
	case '$':
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
	case '*':
		return NewArray(0), frameSize
	case '+':
		return NewSimpleString(text), frameSize
	default:
		return NewSimpleError(fmt.Sprintf("unknown protocol symbol \"%c\"", symbol)), frameSize
	}
}

type SimpleString string

func NewSimpleString(text string) SimpleString {
	return SimpleString(text)
}

func (s SimpleString) IsData() {}

type SimpleError string

func NewSimpleError(s string) SimpleError {
	return SimpleError(s)
}

func (s SimpleError) IsData() {}

type SimpleInteger int64

func NewSimpleInteger(value int64) SimpleInteger {
	return SimpleInteger(value)
}

func (s SimpleInteger) IsData() {}

type BulkString string

func NewBulkString(text string) BulkString {
	return BulkString(text)
}

func (s BulkString) IsData() {}

type Array int

func NewArray(length int) Array {
	return Array(length)
}

func (s Array) IsData() {}
