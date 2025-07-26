package store_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"redis-challenge/internal/store"
	"testing"
)

func TestWritingExpiry(t *testing.T) {

	t.Run("write value into store with no expiry", func(t *testing.T) {
		s := store.New()

		s.WriteWithExpiry("key", "value", store.ExpiryOptionNone, 0)

		confirmKeyHasValue(t, s, "key", "value")
	})

	t.Run("write value into store with expiry time in milliseconds should not exist when expiration time has passed", func(t *testing.T) {
		clock := store.FixedClock{TimeInMilliseconds: 1_000}
		s := store.NewWithCLock(clock.Now)

		s.WriteWithExpiry("key", "value", store.ExpiryOptionExpiryUnixTimeInMilliseconds, 2_000)

		clock.AddSeconds(1).AddMilliseconds(1)

		confirmKeyIsDeleted(t, s, "key")
	})

	t.Run("write value into store with expiry time in milliseconds set into the future", func(t *testing.T) {
		clock := store.FixedClock{TimeInMilliseconds: 1_000}
		s := store.NewWithCLock(clock.Now)

		s.WriteWithExpiry("key", "value", store.ExpiryOptionExpiryUnixTimeInMilliseconds, 2_000)

		confirmKeyHasValue(t, s, "key", "value")
	})

	t.Run("write value into store with expiry time in milliseconds should exist when expiration time is close but not yet passed", func(t *testing.T) {
		clock := store.FixedClock{TimeInMilliseconds: 1_000}
		s := store.NewWithCLock(clock.Now)

		s.WriteWithExpiry("key", "value", store.ExpiryOptionExpiryUnixTimeInMilliseconds, 2_000)

		clock.AddSeconds(1).AddMilliseconds(-1)

		confirmKeyHasValue(t, s, "key", "value")
	})

	t.Run("write value into store with expiry time in seconds should not exist when expiration time has passed", func(t *testing.T) {
		clock := store.FixedClock{TimeInMilliseconds: 1_000}
		s := store.NewWithCLock(clock.Now)

		s.WriteWithExpiry("key", "value", store.ExpiryOptionExpiryUnixTimeInSeconds, 2)

		clock.AddSeconds(1).AddMilliseconds(1)

		confirmKeyIsDeleted(t, s, "key")
	})

	t.Run("write value into store with expiry time in seconds should exist when expiration time has not yet passed", func(t *testing.T) {
		clock := store.FixedClock{TimeInMilliseconds: 1_000}
		s := store.NewWithCLock(clock.Now)

		s.WriteWithExpiry("key", "value", store.ExpiryOptionExpiryUnixTimeInSeconds, 2)

		clock.AddSeconds(1).AddMilliseconds(-1)

		confirmKeyHasValue(t, s, "key", "value")
	})

	t.Run("write value into store with expiry in milliseconds should not exist when expiration time has passed", func(t *testing.T) {
		clock := store.FixedClock{TimeInMilliseconds: 1_000}
		s := store.NewWithCLock(clock.Now)

		s.WriteWithExpiry("key", "value", store.ExpiryOptionExpiryMilliseconds, 1_000)

		clock.AddSeconds(1).AddMilliseconds(1)

		confirmKeyIsDeleted(t, s, "key")
	})

	t.Run("write value into store with expiry in milliseconds should exist when expiration time has not yet passed", func(t *testing.T) {
		clock := store.FixedClock{TimeInMilliseconds: 1_000}
		s := store.NewWithCLock(clock.Now)

		s.WriteWithExpiry("key", "value", store.ExpiryOptionExpiryMilliseconds, 1_000)

		clock.AddSeconds(1).AddMilliseconds(-1)

		confirmKeyHasValue(t, s, "key", "value")
	})

	t.Run("write value into store with expiry in seconds should not exist when expiration time has passed", func(t *testing.T) {
		clock := store.FixedClock{TimeInMilliseconds: 1_000}
		s := store.NewWithCLock(clock.Now)

		s.WriteWithExpiry("key", "value", store.ExpiryOptionExpirySeconds, 1)

		clock.AddSeconds(1).AddMilliseconds(1)

		confirmKeyIsDeleted(t, s, "key")
	})

	t.Run("write value into store with expiry in seconds should exist when expiration time has not yet passed", func(t *testing.T) {
		clock := store.FixedClock{TimeInMilliseconds: 1_000}
		s := store.NewWithCLock(clock.Now)

		s.WriteWithExpiry("key", "value", store.ExpiryOptionExpirySeconds, 1)

		clock.AddSeconds(1).AddMilliseconds(-1)

		confirmKeyHasValue(t, s, "key", "value")
	})

	t.Run("overwrite unexpired value with expiry with KEEPTTL should not exist when expiration time of previous value has passed", func(t *testing.T) {
		clock := store.FixedClock{TimeInMilliseconds: 1_000}
		s := store.NewWithCLock(clock.Now)

		s.WriteWithExpiry("key", "value 1", store.ExpiryOptionExpirySeconds, 1)
		s.WriteWithExpiry("key", "value 2", store.ExpiryOptionExpiryKeepTTL, 0)

		clock.AddSeconds(1).AddMilliseconds(1)

		confirmKeyIsDeleted(t, s, "key")
	})

	t.Run("overwrite unexpired value with expiry with KEEPTTL should exist when expiration time of previous value has not yet passed", func(t *testing.T) {
		clock := store.FixedClock{TimeInMilliseconds: 1_000}
		s := store.NewWithCLock(clock.Now)

		s.WriteWithExpiry("key", "value 1", store.ExpiryOptionExpirySeconds, 1)
		s.WriteWithExpiry("key", "value 2", store.ExpiryOptionExpiryKeepTTL, 0)

		clock.AddSeconds(1).AddMilliseconds(-1)

		confirmKeyHasValue(t, s, "key", "value 2")
	})

	t.Run("write new value with KEEPTTL of should exist", func(t *testing.T) {
		clock := store.FixedClock{TimeInMilliseconds: 1_000}
		s := store.NewWithCLock(clock.Now)

		s.WriteWithExpiry("key", "value 2", store.ExpiryOptionExpiryKeepTTL, 0)

		confirmKeyHasValue(t, s, "key", "value 2")
	})
}

func confirmKeyHasValue(t *testing.T, s store.Store, key string, value string) {
	actualValue, err := s.ReadString(key)

	require.NoError(t, err)
	assert.Equal(t, value, actualValue)

	assert.True(t, s.Exists(key), "key should exist")
}

func confirmKeyIsDeleted(t *testing.T, s store.Store, key string) {
	value, err := s.ReadString(key)

	assert.Equal(t, store.ErrorKeyNotFound, err)
	assert.Equal(t, "", value)

	assert.False(t, s.Exists(key), "key should not exist")
}
