package store_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"redis-challenge/internal/store"
	"testing"
)

func TestTrackingStoreSize(t *testing.T) {

	t.Run("an empty store has a size of 0", func(t *testing.T) {
		s := store.New()
		assert.Equal(t, 0, s.Size())
	})

	t.Run("writing value to empty store increases size to 1", func(t *testing.T) {
		s := store.New()
		s.Write("key", "value", store.ExpiryOptionNone, 0)
		assert.Equal(t, 1, s.Size())
	})

	t.Run("write value, incrementing a second and pushing to a third to an empty list increases size to 3", func(t *testing.T) {
		s := store.New()

		s.Write("key1", "value", store.ExpiryOptionNone, 0)
		_, err := s.Increment("key2", 1)
		require.NoError(t, err)
		_, err = s.LeftPush("key3", []string{"a", "b"})
		require.NoError(t, err)

		assert.Equal(t, 3, s.Size())
	})

	t.Run("deleting a value will decrease the size of the store", func(t *testing.T) {
		s := store.New()

		s.Write("key1", "value", store.ExpiryOptionNone, 0)
		s.Write("key2", "value", store.ExpiryOptionNone, 0)

		sizeBeforeDelete := s.Size()
		s.Delete("key1")
		sizeAfterDelete := s.Size()

		assert.Equal(t, -1, sizeAfterDelete-sizeBeforeDelete, "should be 1 less than before delete")
	})

	t.Run("reading a stored value will not decrease the size of the store", func(t *testing.T) {
		s := store.New()

		s.Write("key1", "value", store.ExpiryOptionNone, 0)
		s.Write("key2", "value", store.ExpiryOptionNone, 0)

		sizeBefore := s.Size()
		_, err := s.ReadString("key1")
		require.NoError(t, err)
		sizeAfter := s.Size()

		assert.Equal(t, 0, sizeAfter-sizeBefore, "reading string value should not decrease the size of the store")
	})

	t.Run("reading an expired value should remove it from the store", func(t *testing.T) {
		clock := store.FixedClock{TimeInMilliseconds: 1_000}
		s := store.NewWithClock(clock.Now)

		s.Write("key1", "value", store.ExpiryOptionExpiryMilliseconds, 1)
		s.Write("key2", "value", store.ExpiryOptionNone, 0)

		clock.AddMilliseconds(2)

		sizeBefore := s.Size()
		_, err := s.ReadString("key1")
		require.Equal(t, store.ErrorKeyNotFound, err)
		sizeAfter := s.Size()

		assert.Equal(t, -1, sizeAfter-sizeBefore, "reading string value should not decrease the size of the store")
	})
}
