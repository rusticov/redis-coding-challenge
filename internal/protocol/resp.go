package protocol

import (
	"bytes"
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

	switch bs[0] {
	case '-':
		return NewError(text), delimiterIndex + 2
	case ':':
		value, _ := strconv.ParseInt(text, 10, 64)
		return NewSimpleInteger(value), delimiterIndex + 2
	default:
		return NewSimpleString(text), delimiterIndex + 2
	}
}

type SimpleString string

func NewSimpleString(s string) SimpleString {
	return SimpleString(s)
}

func (s SimpleString) IsData() {}

type Error string

func NewError(s string) Error {
	return Error(s)
}

func (s Error) IsData() {}

type SimpleInteger int64

func NewSimpleInteger(value int64) SimpleInteger {
	return SimpleInteger(value)
}

func (s SimpleInteger) IsData() {}
