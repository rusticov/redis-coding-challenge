package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Call represents an interface for defining protocol interactions with requests and responses for a server or system.
// It allows tests described in protocol or data forms to both be run against a running Redis server or server clone.
type Call interface {
	Request() string
	IsResponseExpected() bool
	ConfirmResponse(t testing.TB, response string)
}

func NewCallWithProtocol(request string, expectedResponse string) ProtocolStringCall {
	return ProtocolStringCall{
		request:          request,
		expectedResponse: expectedResponse,
	}
}

func NewCallWithProtocolAndPartialResponse(request string, expectedResponse string) ProtocolStringCall {
	return ProtocolStringCall{
		request:                 request,
		expectedResponse:        expectedResponse,
		realResponseMayBeLonger: true,
	}
}

func NewCallWithProtocolWithoutResponse(request string) ProtocolStringCall {
	return ProtocolStringCall{
		request: request,
	}
}

type ProtocolStringCall struct {
	request                 string
	expectedResponse        string
	realResponseMayBeLonger bool
}

func (p ProtocolStringCall) Request() string {
	return p.request
}

func (p ProtocolStringCall) IsResponseExpected() bool {
	return p.expectedResponse != ""
}

func (p ProtocolStringCall) ConfirmResponse(t testing.TB, response string) {
	if p.realResponseMayBeLonger {
		assert.Contains(t, response, p.expectedResponse)
	} else {
		assert.Equal(t, p.expectedResponse, response)
	}
}
