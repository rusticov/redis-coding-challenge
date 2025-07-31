package store_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"redis-challenge/internal/store"
	"testing"
)

func TestPushingListsIntoStore(t *testing.T) {

	t.Run("left push to new key should create list", func(t *testing.T) {
		// Given empty store
		s := store.New()

		// When a list is left pushed
		count, err := s.LeftPush("key", []string{"a", "b"})

		// Then push succeeds and the new list can be read
		require.NoError(t, err)
		assert.Equal(t, int64(2), count)

		listRange, err := s.ReadListRange("key", 0, -1)
		require.NoError(t, err, "list range can be read")
		assert.Equal(t, []string{"b", "a"}, listRange, "list range should be read with value in reverse order")
	})

	t.Run("left push to key with list should push to start of the list", func(t *testing.T) {
		// Given store has a key with a list
		s := store.New()
		_, err := s.LeftPush("key", []string{"a", "b"})
		require.NoError(t, err)

		// When a list is left pushed to that key
		count, err := s.LeftPush("key", []string{"c", "d", "e"})

		// Then push succeeds and the new list can be read
		require.NoError(t, err)
		assert.Equal(t, int64(2+3), count, "count of total items in list")

		listRange, err := s.ReadListRange("key", 0, -1)
		require.NoError(t, err, "list range can be read")
		assert.Equal(t, []string{"e", "d", "c", "b", "a"}, listRange, "list range should be read with value in reverse order")
	})

	t.Run("left push against key with string value should not create list", func(t *testing.T) {
		// Given empty store
		s := store.New()

		s.Write("key", "value", store.ExpiryOptionNone, 0)

		// When a list is left pushed
		_, err := s.LeftPush("key", []string{"a", "b"})

		// Then the push should fail
		assert.Equal(t, store.ErrorWrongOperationType, err)
	})

	t.Run("left push an expired value should set a new value with no expiry", func(t *testing.T) {
		// Given store with an expired string value
		clock := &store.FixedClock{TimeInMilliseconds: 1_000}
		s := store.NewWithClock(clock)

		s.Write("key", "value", store.ExpiryOptionExpiryMilliseconds, 1_000)
		clock.AddSeconds(1).AddMilliseconds(1)

		// When a list is left pushed
		count, err := s.LeftPush("key", []string{"a", "b"})

		// Then push succeeds and the new list can be read
		require.NoError(t, err)
		assert.Equal(t, int64(2), count)

		listRange, err := s.ReadListRange("key", 0, -1)
		require.NoError(t, err, "list range can be read")
		assert.Equal(t, []string{"b", "a"}, listRange, "list range should be read")
	})

	t.Run("right push to new key should create list", func(t *testing.T) {
		// Given empty store
		s := store.New()

		// When a list is right pushed
		count, err := s.RightPush("key", []string{"a", "b"})

		// Then push succeeds and the new list can be read
		require.NoError(t, err)
		assert.Equal(t, int64(2), count)

		listRange, err := s.ReadListRange("key", 0, -1)
		require.NoError(t, err, "list range can be read")
		assert.Equal(t, []string{"a", "b"}, listRange, "list range should be read with values in order")
	})

	t.Run("right push to key with list should push to start of the list", func(t *testing.T) {
		// Given store has a key with a list
		s := store.New()
		_, err := s.RightPush("key", []string{"a", "b"})
		require.NoError(t, err)

		// When a list is right pushed to that key
		count, err := s.RightPush("key", []string{"c", "d", "e"})

		// Then push succeeds and the new list can be read
		require.NoError(t, err)
		assert.Equal(t, int64(2+3), count, "count of total items in list")

		listRange, err := s.ReadListRange("key", 0, -1)
		require.NoError(t, err, "list range can be read")
		assert.Equal(t, []string{"a", "b", "c", "d", "e"}, listRange, "list range should be read with value in order")
	})

	t.Run("right push against key with string value should not create list", func(t *testing.T) {
		// Given empty store
		s := store.New()

		s.Write("key", "value", store.ExpiryOptionNone, 0)

		// When a list is right pushed
		_, err := s.RightPush("key", []string{"a", "b"})

		// Then the push should fail
		assert.Equal(t, store.ErrorWrongOperationType, err)
	})

	t.Run("right push an expired value should set a new value with no expiry", func(t *testing.T) {
		// Given store with an expired string value
		clock := &store.FixedClock{TimeInMilliseconds: 1_000}
		s := store.NewWithClock(clock)

		s.Write("key", "value", store.ExpiryOptionExpiryMilliseconds, 1_000)
		clock.AddSeconds(1).AddMilliseconds(1)

		// When a list is right pushed
		count, err := s.RightPush("key", []string{"a", "b"})

		// Then push succeeds and the new list can be read
		require.NoError(t, err)
		assert.Equal(t, int64(2), count)

		listRange, err := s.ReadListRange("key", 0, -1)
		require.NoError(t, err, "list range can be read")
		assert.Equal(t, []string{"a", "b"}, listRange, "list range should be read")
	})

	t.Run("reading a string value against a list should error", func(t *testing.T) {
		// Given empty store
		s := store.New()
		_, err := s.LeftPush("key", []string{"a", "b"})
		require.NoError(t, err)

		// When reading the value as a string
		_, err = s.ReadString("key")

		// Then an error is returned
		assert.Equal(t, store.ErrorWrongOperationType, err)
	})

	t.Run("incrementing against a list should error", func(t *testing.T) {
		// Given empty store
		s := store.New()
		_, err := s.LeftPush("key", []string{"a", "b"})
		require.NoError(t, err)

		// When reading the value as a string
		_, err = s.Increment("key", 1)

		// Then an error is returned
		assert.Equal(t, store.ErrorWrongOperationType, err)
	})
}
