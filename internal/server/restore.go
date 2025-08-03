package server

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"redis-challenge/internal/command"
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
)

type restorer struct {
	store store.Store
}

func (h restorer) RestoreFromLog(reader io.Reader) error {
	var buffer bytes.Buffer

	var totalReadByteCount int
	readBuffer := make([]byte, 32768)

readMore:
	for {
		bytesRead, err := reader.Read(readBuffer)
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return fmt.Errorf("failed to read request: %w", err)
		}
		totalReadByteCount += bytesRead
		slog.Info("read bytes", "bytes", bytesRead, "total", totalReadByteCount)

		buffer.Write(readBuffer[:bytesRead])

		readBytes := buffer.Bytes()

		offset := 0
		for {
			protocolData, requestByteCount := protocol.ReadFrame(readBytes[offset:])
			if requestByteCount == 0 {
				remainingBytes := readBytes[offset:]
				buffer.Reset()
				buffer.Write(remainingBytes)
				continue readMore
			}

			err = h.executeCommand(protocolData, readBytes[offset:offset+requestByteCount])
			if err != nil {
				return err
			}

			offset += requestByteCount
		}
	}
}

func (h restorer) executeCommand(protocolData protocol.Data, requestBytes []byte) error {
	parsedCommand, commandError := command.Validate(protocolData)

	switch {
	case commandError != nil:
		return fmt.Errorf("failed to parse request from log: %v %s", commandError, string(requestBytes))
	case parsedCommand == nil:
		return fmt.Errorf("request from log is not a command: %v %s", commandError, string(requestBytes))
	default:
		_, err := parsedCommand.Execute(h.store)
		if err != nil {
			return err
		}
		return nil
	}
}
