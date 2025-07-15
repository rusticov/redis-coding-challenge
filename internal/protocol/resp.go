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
	return NewSimpleString(string(bs[1:delimiterIndex])), delimiterIndex + 2
}

type SimpleString struct{}

func NewSimpleString(s string) Data {
	return nil
}

func (s SimpleString) IsData() {}
