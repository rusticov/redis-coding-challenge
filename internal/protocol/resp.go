package protocol

import (
	"bytes"
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

	switch bs[0] {
	case '-':
		return NewError(string(bs[1:delimiterIndex])), delimiterIndex + 2
	default:
		return NewSimpleString(string(bs[1:delimiterIndex])), delimiterIndex + 2
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
