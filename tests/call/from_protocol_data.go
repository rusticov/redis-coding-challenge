package call

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"redis-challenge/internal/protocol"
	"testing"
)

func NewFromData(request []protocol.Data, expectedResponse protocol.Data) DataCall {
	return DataCall{
		request:          request,
		expectedResponse: expectedResponse,
	}
}

func NewFromDataWithPartialError(request []protocol.Data, expectedError string) DataCall {
	return DataCall{
		request:              request,
		expectedPartialError: expectedError,
	}
}

func NewFromDataWithoutError(request []protocol.Data) DataCall {
	return DataCall{
		request:          request,
		callIsNotAnError: true,
	}
}

type DataCall struct {
	request              []protocol.Data
	expectedResponse     protocol.Data
	expectedPartialError string
	callIsNotAnError     bool
}

func (c DataCall) RequestData() protocol.Data {
	return protocol.Array{Data: c.request}
}

func (c DataCall) Request() string {
	var buffer bytes.Buffer
	err := protocol.WriteData(&buffer, c.RequestData())

	if err != nil {
		return err.Error()
	}
	return buffer.String()
}

func (c DataCall) IsResponseExpected() bool {
	return c.expectedResponse != nil
}

func (c DataCall) ConfirmResponse(t testing.TB, response string) {
	switch {
	case c.callIsNotAnError:
		assert.NotEqual(t, "-", response[0:1], "response should not be an error")
	case c.expectedPartialError != "":
		assert.Contains(t, response, "-"+c.expectedPartialError)
	default:
		actualResponse, _ := protocol.ReadFrame([]byte(response))
		assert.Equal(t, c.expectedResponse, actualResponse)
	}
}

func (c DataCall) ConfirmValidation(t testing.TB, validationError protocol.Data) {
	if c.callIsNotAnError {
		assert.Nil(t, validationError, "response should not be an error")
		return
	}

	expectedError := c.expectedResponse
	if c.expectedPartialError != "" {
		expectedError = protocol.NewSimpleError(c.expectedPartialError)
	}

	assert.Equal(t, expectedError, validationError, "response should be an error")
}
