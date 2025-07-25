package call

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func NewFromProtocol(request string, expectedResponse string) Call {
	return protocolStringCall{
		request:          request,
		expectedResponse: expectedResponse,
	}
}

func NewFromProtocolWithPartialResponse(request string, expectedResponse string) Call {
	return protocolStringCall{
		request:                 request,
		expectedResponse:        expectedResponse,
		realResponseMayBeLonger: true,
	}
}

func NewFromProtocolWithoutResponse(request string) Call {
	return protocolStringCall{
		request: request,
	}
}

type protocolStringCall struct {
	request                 string
	expectedResponse        string
	realResponseMayBeLonger bool
}

func (p protocolStringCall) Request() string {
	return p.request
}

func (p protocolStringCall) Delay() time.Duration {
	return 0
}

func (p protocolStringCall) IsResponseExpected() bool {
	return p.expectedResponse != ""
}

func (p protocolStringCall) ConfirmResponse(t testing.TB, response string) {
	if p.realResponseMayBeLonger {
		assert.Contains(t, response, p.expectedResponse)
	} else {
		assert.Equal(t, p.expectedResponse, response)
	}
}

func (p protocolStringCall) IsPossiblePartialResponse(response string) bool {
	return response != p.expectedResponse && strings.HasPrefix(p.expectedResponse, response)
}
