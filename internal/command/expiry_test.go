package command_test

import (
	"github.com/stretchr/testify/assert"
	"redis-challenge/internal/command"
	"redis-challenge/internal/store"
	"testing"
)

func TestExpiry(t *testing.T) {

	clock := store.FixedClock{TimeInMilliseconds: 10_456}

	t.Run("with no expiry", func(t *testing.T) {
		option, timestamp := command.ExpiryTimestamp(clock.Now, store.ExpiryOptionNone, 0)

		assert.Equal(t, store.ExpiryOptionNone, option)
		assert.Equal(t, int64(0), timestamp)
	})

	t.Run("with keep TTL", func(t *testing.T) {
		option, timestamp := command.ExpiryTimestamp(clock.Now, store.ExpiryOptionExpiryKeepTTL, 0)

		assert.Equal(t, store.ExpiryOptionExpiryKeepTTL, option)
		assert.Equal(t, int64(0), timestamp)
	})

	t.Run("with expiry in unix milliseconds should return timestamp as given", func(t *testing.T) {
		option, timestamp := command.ExpiryTimestamp(clock.Now, store.ExpiryOptionExpiryUnixTimeInMilliseconds, 42_123)

		assert.Equal(t, store.ExpiryOptionExpiryUnixTimeInMilliseconds, option)
		assert.Equal(t, int64(42_123), timestamp)
	})

	t.Run("with expiry in unix seconds should return timestamp in milliseconds", func(t *testing.T) {
		option, timestamp := command.ExpiryTimestamp(clock.Now, store.ExpiryOptionExpiryUnixTimeInSeconds, 42)

		assert.Equal(t, store.ExpiryOptionExpiryUnixTimeInMilliseconds, option)
		assert.Equal(t, int64(42_000), timestamp)
	})

	t.Run("with expiry of milliseconds in the future should return timestamp in milliseconds", func(t *testing.T) {
		option, timestamp := command.ExpiryTimestamp(clock.Now, store.ExpiryOptionExpiryMilliseconds, 42_123)

		assert.Equal(t, store.ExpiryOptionExpiryUnixTimeInMilliseconds, option)
		assert.Equal(t, int64(52_579), timestamp)
	})

	t.Run("with expiry of seconds in the future should return timestamp in milliseconds", func(t *testing.T) {
		option, timestamp := command.ExpiryTimestamp(clock.Now, store.ExpiryOptionExpirySeconds, 42)

		assert.Equal(t, store.ExpiryOptionExpiryUnixTimeInMilliseconds, option)
		assert.Equal(t, int64(52_456), timestamp)
	})
}
