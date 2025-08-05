package protocol

import "redis-challenge/internal/list"

type Data interface {
	Symbol() DataTypeSymbol
}

type SimpleString string

func NewSimpleString(text string) SimpleString {
	return SimpleString(text)
}

func (s SimpleString) Symbol() DataTypeSymbol {
	return SimpleStringSymbol
}

type SimpleError string

func NewSimpleError(s string) SimpleError {
	return SimpleError(s)
}

func (s SimpleError) Symbol() DataTypeSymbol {
	return SimpleErrorSymbol
}

type SimpleInteger int64

func NewSimpleInteger(value int64) SimpleInteger {
	return SimpleInteger(value)
}

func (s SimpleInteger) Symbol() DataTypeSymbol {
	return SimpleIntegerSymbol
}

type BulkString string

func NewBulkString(text string) BulkString {
	return BulkString(text)
}

func (s BulkString) Symbol() DataTypeSymbol {
	return BulkStringSymbol
}

type Array struct {
	Data []Data
}

func NewArray(data []Data) Array {
	return Array{Data: data}
}

func (s Array) Symbol() DataTypeSymbol {
	return ArraySymbol
}

type DoubleEndedList struct {
	Data list.DoubleEndedList
}

func (s DoubleEndedList) Symbol() DataTypeSymbol {
	return ArraySymbol
}

func NewDoubleEndedList(list list.DoubleEndedList) Data {
	return DoubleEndedList{Data: list}
}
