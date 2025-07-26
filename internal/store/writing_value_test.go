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

	t.Run("write value into store with expiry seconds set into the future", func(t *testing.T) {
		clock := store.FixedClock{TimeInMilliseconds: 1000}
		s := store.NewWithCLock(clock.Now)

		s.WriteWithExpiry("key", "value", store.ExpiryOptionExpiryUnixTimeInSeconds, 1001)

		confirmKeyHasValue(t, s, "key", "value")
	})

	t.Run("write value into store with expiry seconds should not exist when expiration time has passed", func(t *testing.T) {
		clock := store.FixedClock{TimeInMilliseconds: 1000}
		s := store.NewWithCLock(clock.Now)

		s.WriteWithExpiry("key", "value", store.ExpiryOptionExpiryUnixTimeInSeconds, 1001)

		clock.AddSeconds(1).AddMilliseconds(1)

		confirmKeyIsDeleted(t, s, "key")
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
