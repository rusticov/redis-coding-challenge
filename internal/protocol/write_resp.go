package protocol

import (
	"fmt"
	"io"
	"strconv"
)

func WriteData(out io.Writer, data Data) error {
	var text string
	switch d := data.(type) {
	case nil:
		text = "$-1\r\n"
	case SimpleString:
		return writeString(out, SimpleStringSymbol, string(d))
	case SimpleError:
		return writeString(out, SimpleErrorSymbol, string(d))
	case SimpleInteger:
		return writeNumber(out, SimpleIntegerSymbol, int64(d))
	case BulkString:
		return writeBulkString(out, d)
	case Array:
		return writeArray(out, d)
	case DoubleEndedList:
		return writeDoubleEndedList(out, d)
	default:
		text = fmt.Sprintf("-ERR unknown data type\r\n")
	}

	_, err := out.Write([]byte(text))
	return err
}

func writeString(out io.Writer, symbol DataTypeSymbol, text string) error {
	if _, err := out.Write([]byte{byte(symbol)}); err != nil {
		return err
	}

	if _, err := out.Write([]byte(text)); err != nil {
		return err
	}

	if _, err := out.Write([]byte("\r\n")); err != nil {
		return err
	}
	return nil
}

func writeNumber(out io.Writer, symbol DataTypeSymbol, number int64) error {
	return writeString(out, symbol, strconv.FormatInt(number, 10))
}

func writeBulkString(out io.Writer, data BulkString) error {
	if err := writeNumber(out, BulkStringSymbol, int64(len(data))); err != nil {
		return err
	}

	if _, err := out.Write([]byte(data)); err != nil {
		return err
	}
	if _, err := out.Write([]byte("\r\n")); err != nil {
		return err
	}
	return nil
}

func writeArray(out io.Writer, d Array) error {
	if err := writeNumber(out, ArraySymbol, int64(len(d.Data))); err != nil {
		return err
	}

	for _, item := range d.Data {
		if err := WriteData(out, item); err != nil {
			return err
		}
	}

	return nil
}

func writeDoubleEndedList(out io.Writer, d DoubleEndedList) error {
	if err := writeNumber(out, ArraySymbol, int64(d.Data.Len())); err != nil {
		return err
	}

	for _, item := range d.Data.Range() {
		if err := WriteData(out, BulkString(item)); err != nil {
			return err
		}
	}

	return nil
}
