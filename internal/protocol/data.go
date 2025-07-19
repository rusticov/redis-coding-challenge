package protocol

type Data interface {
	IsData()
	Symbol() rune
}

type SimpleString string

func NewSimpleString(text string) SimpleString {
	return SimpleString(text)
}

func (s SimpleString) IsData() {}

func (s SimpleString) Symbol() rune {
	return '+'
}

type SimpleError string

func NewSimpleError(s string) SimpleError {
	return SimpleError(s)
}

func (s SimpleError) IsData() {}

func (s SimpleError) Symbol() rune {
	return '-'
}

type SimpleInteger int64

func NewSimpleInteger(value int64) SimpleInteger {
	return SimpleInteger(value)
}

func (s SimpleInteger) IsData() {}

func (s SimpleInteger) Symbol() rune {
	return ':'
}

type BulkString string

func NewBulkString(text string) BulkString {
	return BulkString(text)
}

func (s BulkString) IsData() {}

func (s BulkString) Symbol() rune {
	return '$'
}

type Array struct {
	Data []Data
}

func NewArray(data []Data) Array {
	return Array{Data: data}
}

func (s Array) IsData() {}

func (s Array) Symbol() rune {
	return '*'
}
