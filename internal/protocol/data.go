package protocol

type Data interface {
	IsData()
	Symbol() DataTypeSymbol
}

type SimpleString string

func NewSimpleString(text string) SimpleString {
	return SimpleString(text)
}

func (s SimpleString) IsData() {}

func (s SimpleString) Symbol() DataTypeSymbol {
	return SimpleStringSymbol
}

type SimpleError string

func NewSimpleError(s string) SimpleError {
	return SimpleError(s)
}

func (s SimpleError) IsData() {}

func (s SimpleError) Symbol() DataTypeSymbol {
	return SimpleErrorSymbol
}

type SimpleInteger int64

func NewSimpleInteger(value int64) SimpleInteger {
	return SimpleInteger(value)
}

func (s SimpleInteger) IsData() {}

func (s SimpleInteger) Symbol() DataTypeSymbol {
	return SimpleIntegerSymbol
}

type BulkString string

func NewBulkString(text string) BulkString {
	return BulkString(text)
}

func (s BulkString) IsData() {}

func (s BulkString) Symbol() DataTypeSymbol {
	return BulkStringSymbol
}

type Array struct {
	Data []Data
}

func NewArray(data []Data) Array {
	return Array{Data: data}
}

func (s Array) IsData() {}

func (s Array) Symbol() DataTypeSymbol {
	return ArraySymbol
}
