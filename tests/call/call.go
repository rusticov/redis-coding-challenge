package call

import (
	"testing"
)

// Call represents an interface for defining protocol interactions with requests and responses for a server or system.
// It allows tests described in protocol or data forms to both be run against a running Redis server or server clone.
type Call interface {
	Request() string
	IsResponseExpected() bool
	ConfirmResponse(t testing.TB, response string)
	IsPossiblePartialResponse(response string) bool
}
