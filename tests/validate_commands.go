package tests

import (
	"redis-challenge/internal/command"
	"redis-challenge/tests/call"
	"testing"
)

type SelectTestCaseDriver string

const (
	SelectTestCaseDriverRedisServer SelectTestCaseDriver = "redis-server-driver"
	SelectTestCaseDriverRedisClone  SelectTestCaseDriver = "redis-clone-driver"
)

func ValidateCommands(t testing.TB, calls []call.DataCall, driverChoice SelectTestCaseDriver) {
	switch driverChoice {
	case SelectTestCaseDriverRedisServer:
		DriveProtocolAgainstServer(t, calls, UseRealRedisServer)
	case SelectTestCaseDriverRedisClone:
		DriveProtocolAgainstServer(t, calls, UseChallengeServer)
	default:
		validateAgainstCommandValidator(t, calls)
	}
}

func validateAgainstCommandValidator(t testing.TB, calls []call.DataCall) {
	for _, c := range calls {
		_, errorData := command.Validate(c.RequestData())

		c.ConfirmValidation(t, errorData)
	}
}
