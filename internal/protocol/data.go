package protocol

type Data interface {
	IsData()
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

type Array struct {
	Data []Data
}

func NewArray(data []Data) Array {
	return Array{Data: data}
}

func (s Array) IsData() {}
