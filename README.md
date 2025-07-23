# Redis Challenge

This project is an implementation of a Redis server for the "Coding Challenges"
from [https://codingchallenges.fyi/challenges/challenge-wc/](https://codingchallenges.fyi/challenges/challenge-wc/).
It is being developed as part of the accompanying course.

## Currently Supported Commands

The following Redis commands are currently implemented:

* PING
* ECHO
* GET
* SET
* DEL
* EXISTS
* INCR
* DECR
* LPUSH
* RPUSH
* LRANGE

## Build and Running Server

The server currently recognizes no arguments.  It runs against a random port to
avoid clashing with Redis whilst in development.

```bash
# Build the project
go build -v ./...
```

## Running Tests

```bash
# Run all tests (uses challenge clone implementation by default)
go test -v ./...
```

## Requirements

- Go 1.24 or later
- Redis server (for comparative testing only)

## Project Structure

- `internal/command/` - Command implementations (PING, ECHO, GET, SET, etc)
- `internal/protocol/` - Redis protocol parsing and serialization
- `internal/server/` - Server implementation
- `internal/store/` - Key-value store implementation
- `tests/` - Test utilities and high-level test cases

## Testing Against Real Redis Server

This project includes the ability to test the implementation against a real running Redis server for comparison.
You can switch between testing your implementation and a real Redis server using the driver/variant options.

### Setting Up Testing

To test against a real Redis server, you need:

1. A Redis server running on `localhost:6379` (default Redis port)
2. Use the `ServerVariant` parameter in your tests

### Testing Examples

The high-level tests are set up so that they can be run against a running instance of Redis.
This means that the tests can be verified as accurate before attempting to produce the
cloned implementation.  Enums are used to switch between using Redis and code.

A typical flow will mean that a test is first verified as accurate by running against Redis, before
then using that test to generate implementation in the clone.

#### Basic Test Structure

```go
import (
    "redis-challenge/tests"
    "redis-challenge/tests/call"
)

// Test against the challenge implementation (default)
tests.DriveProtocolAgainstServer(t, testCalls)

// Test against a real Redis server
tests.DriveProtocolAgainstServer(t, testCalls, tests.UseRealRedisServer)
```

#### Using Test Cases with Driver Choice

```go
type testCase struct {
    calls        []call.DataCall
    driverChoice tests.ServerVariant
}

testCases := map[string]testCase{
    "test with challenge server": {
        calls: myTestCalls,
        driverChoice: tests.UseChallengeServer, // or omit for default
    },
    "test with real redis": {
        calls: myTestCalls,
        driverChoice: tests.UseRealRedisServer,
    },
}

for name, tc := range testCases {
    t.Run(name, func(t *testing.T) {
        tests.DriveProtocolAgainstServer(t, tc.calls, tc.driverChoice)
    })
}
```
