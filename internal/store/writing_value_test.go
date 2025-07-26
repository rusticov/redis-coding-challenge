package store_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"redis-challenge/internal/store"
	"testing"
)

func TestWritingExpiry(t *testing.T) {

	t.Run("read value written into store with no expiry", func(t *testing.T) {
		s := store.New()

		s.WriteWithExpiry("key", "value", store.ExpiryOptionNone, 0)
		value, err := s.ReadString("key")

		require.NoError(t, err)
		assert.Equal(t, "value", value)
	})

	t.Run("read value written into store with expiry seconds set into the future", func(t *testing.T) {
		clock := store.FixedClock{TimeInMilliseconds: 1000}
		s := store.NewWithCLock(clock.Now)

		s.WriteWithExpiry("key", "value", store.ExpiryOptionExpiryUnixTimeInSeconds, 1001)
		value, err := s.ReadString("key")

		require.NoError(t, err)
		assert.Equal(t, "value", value)
	})

	t.Run("read value written into store with expiry seconds should not be found when expiration time has passed", func(t *testing.T) {
		clock := store.FixedClock{TimeInMilliseconds: 1000}
		s := store.NewWithCLock(clock.Now)

		s.WriteWithExpiry("key", "value", store.ExpiryOptionExpiryUnixTimeInSeconds, 1001)

		clock.AddSeconds(1).AddMilliseconds(1)
		value, err := s.ReadString("key")

		assert.Equal(t, store.ErrorKeyNotFound, err, "should not find key")
		assert.Equal(t, "", value, "should not find value")
	})
}
