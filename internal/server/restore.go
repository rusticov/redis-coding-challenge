package server

import (
	"bytes"
	"fmt"
	"io"
	"redis-challenge/internal/command"
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
)

type restorer struct {
	store store.Store
}

func (h restorer) RestoreFromLog(reader io.Reader) error {
	var buffer bytes.Buffer

	readBuffer := make([]byte, 1024)

readMore:
	for {
		bytesRead, err := reader.Read(readBuffer)
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return fmt.Errorf("failed to read request: %w", err)
		}
		buffer.Write(readBuffer[:bytesRead])

		for {
			protocolData, requestByteCount := protocol.ReadFrame(readBuffer)
			if requestByteCount == 0 {
				continue readMore
			}

			err = h.executeCommand(protocolData, buffer)
			if err != nil {
				return err
			}

			copy(readBuffer, readBuffer[requestByteCount:])
		}
	}
}

func (h restorer) executeCommand(protocolData protocol.Data, buffer bytes.Buffer) error {
	parsedCommand, commandError := command.Validate(protocolData)

	switch {
	case commandError != nil:
		return fmt.Errorf("failed to parse request from log: %v %s", commandError, buffer.String())
	case parsedCommand == nil:
		return fmt.Errorf("request from log is not a command: %v %s", commandError, buffer.String())
	default:
		_, err := parsedCommand.Execute(h.store)
		if err != nil {
			return err
		}
		return nil
	}
}
