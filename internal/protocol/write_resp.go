package protocol

import (
	"fmt"
	"io"
)

func WriteData(out io.Writer, data Data) error {
	var text string
	switch d := data.(type) {
	case nil:
		text = "$-1\r\n"
	case SimpleString:
		text = fmt.Sprintf("+%s\r\n", d)
	case SimpleError:
		text = fmt.Sprintf("-%s\r\n", d)
	case SimpleInteger:
		text = fmt.Sprintf(":%d\r\n", d)
	case BulkString:
		text = fmt.Sprintf("$%d\r\n%s\r\n", len(d), d)
	case Array:
		text = fmt.Sprintf("*%d\r\n", len(d.Data))
		_, err := out.Write([]byte(text))
		if err != nil {
			return err
		}
		for _, item := range d.Data {
			if err := WriteData(out, item); err != nil {
				return err
			}
		}
		return nil
	default:
		text = fmt.Sprintf("-ERR unknown data type\r\n")
	}

	_, err := out.Write([]byte(text))
	return err
}
