package call

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"redis-challenge/internal/protocol"
	"strings"
	"testing"
	"time"
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
	delay                time.Duration
}

func (c DataCall) WithDelay(delay time.Duration) DataCall {
	c.delay = delay
	return c
}

func (c DataCall) Delay() time.Duration {
	return c.delay
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
	return true
}

func (c DataCall) ConfirmResponse(t testing.TB, response string) {
	switch {
	case c.callIsNotAnError:
		assert.NotEqual(t, "-", response[0:1],
			"response should not be an error to the request %s", c.Request())
	case c.expectedPartialError != "":
		assert.Contains(t, response, "-"+c.expectedPartialError,
			"partial error response to the request %s", c.Request())
	default:
		actualResponse, _ := protocol.ReadFrame([]byte(response))
		assert.Equal(t, c.expectedResponse, actualResponse,
			"error response to the request %s", c.Request())
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

func (c DataCall) IsPossiblePartialResponse(response string) bool {
	expectedResponse := c.expectedResponseAsText()
	return response != expectedResponse && strings.HasPrefix(expectedResponse, response)
}

func (c DataCall) expectedResponseAsText() string {
	var buffer bytes.Buffer
	err := protocol.WriteData(&buffer, c.expectedResponse)

	if err != nil {
		return err.Error()
	}
	return buffer.String()
}
