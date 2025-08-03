package server

import (
	"bytes"
	"io"
	"log/slog"
	"net"
	"redis-challenge/internal/command"
	"redis-challenge/internal/protocol"
)

type connectionHandler struct {
	executor command.Executor
}

func (h connectionHandler) HandleConnection(connection net.Conn) {
	defer func() {
		err := connection.Close()
		if err != nil {
			slog.Error("failed to close connection", "error", err)
		}
	}()

	var buffer bytes.Buffer

	readBuffer := make([]byte, 1024)
	for {
		bytesRead, err := connection.Read(readBuffer)
		if err != nil {
			if err == io.EOF {
				return
			}

			slog.Error("failed to read request", "error", err)
			return
		}
		buffer.Write(readBuffer[:bytesRead])

		protocolData, requestByteCount := protocol.ReadFrame(buffer.Bytes())
		if requestByteCount == 0 {
			continue
		}
		response := h.executeCommand(protocolData, buffer)

		outBuffer := bytes.NewBuffer(nil)
		err = protocol.WriteData(outBuffer, response)
		if err != nil {
			slog.Error("failed to write parse response error", "error", err, "request", buffer.String())
		}

		_, err = connection.Write(outBuffer.Bytes())
		if err != nil {
			slog.Error("failed to send response", "error", err, "request", buffer.String())
		}

		copy(readBuffer, buffer.Bytes()[requestByteCount:])
		buffer.Reset()
	}
}

func (h connectionHandler) executeCommand(protocolData protocol.Data, buffer bytes.Buffer) protocol.Data {
	parsedCommand, commandError := command.Validate(protocolData)

	switch {
	case commandError != nil:
		slog.Error("failed to parse request", "error", commandError, "request", buffer.String())
		return commandError
	case parsedCommand == nil:
		slog.Error("expect a command if there is no error data on parsing", "error", commandError, "request", buffer.String())
		return protocol.NewSimpleError("ERR protocol error")
	default:
		responseReceiver := make(chan protocol.Data)
		errorReceiver := make(chan error)

		h.executor.Execute(buffer.Bytes(), parsedCommand, responseReceiver, errorReceiver)

		select {
		case err := <-errorReceiver:
			slog.Error("failed to execute request", "error", err, "request", buffer.String())
			return protocol.NewSimpleError("ERR protocol error")

		case response := <-responseReceiver:
			return response
		}
	}
}
