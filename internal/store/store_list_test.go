package store_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"redis-challenge/internal/store"
	"testing"
)

func TestStoreList(t *testing.T) {

	t.Run("left push to empty key", func(t *testing.T) {
		s := store.New()
		count, err := s.LeftPush("key", []string{"a", "b"})
		require.NoError(t, err)

		assert.Equal(t, int64(2), count, "count of total items in list")
	})

	t.Run("left push to empty key", func(t *testing.T) {
		s := store.New()
		_, err := s.LeftPush("key", []string{"a", "b"})
		require.NoError(t, err)

		values, err := s.ReadListRange("key", 0, 1)
		require.NoError(t, err)

		assert.Equal(t, []string{"b", "a"}, values)
	})

	t.Run("left push to key with a list", func(t *testing.T) {
		s := store.New()

		_, err := s.LeftPush("key", []string{"a", "b"})
		require.NoError(t, err)

		count, err := s.LeftPush("key", []string{"c", "d", "e"})
		require.NoError(t, err)

		assert.Equal(t, int64(5), count, "count of total items in list")
	})

	t.Run("right push to empty key", func(t *testing.T) {
		s := store.New()
		count, err := s.RightPush("key", []string{"a", "b"})
		require.NoError(t, err)

		assert.Equal(t, int64(2), count, "count of total items in list")
	})

	t.Run("right push to empty key", func(t *testing.T) {
		s := store.New()
		_, err := s.RightPush("key", []string{"a", "b"})
		require.NoError(t, err)

		values, err := s.ReadListRange("key", 0, 1)
		require.NoError(t, err)

		assert.Equal(t, []string{"a", "b"}, values)
	})

	t.Run("right push to key with a list", func(t *testing.T) {
		s := store.New()

		_, err := s.RightPush("key", []string{"a", "b"})
		require.NoError(t, err)

		count, err := s.RightPush("key", []string{"c", "d", "e"})
		require.NoError(t, err)

		assert.Equal(t, int64(5), count, "count of total items in list")
	})

	t.Run("right push to key with a list", func(t *testing.T) {
		s := store.New()

		_, err := s.RightPush("key", []string{"a", "b"})
		require.NoError(t, err)

		count, err := s.RightPush("key", []string{"c", "d", "e"})
		require.NoError(t, err)

		assert.Equal(t, int64(5), count, "count of total items in list")
	})

	t.Run("right push to key with a list should to the end in order", func(t *testing.T) {
		s := store.New()

		_, err := s.RightPush("key", []string{"a", "b"})
		require.NoError(t, err)

		_, err = s.RightPush("key", []string{"c", "d", "e"})
		require.NoError(t, err)

		values, err := s.ReadListRange("key", 0, 10)
		require.NoError(t, err)

		assert.Equal(t, []string{"a", "b", "c", "d", "e"}, values)
	})

	t.Run("read sub-range of list", func(t *testing.T) {
		s := store.New()

		_, err := s.RightPush("key", []string{"a", "b"})
		require.NoError(t, err)

		_, err = s.RightPush("key", []string{"c", "d", "e"})
		require.NoError(t, err)

		values, err := s.ReadListRange("key", 2, -2)
		require.NoError(t, err)

		assert.Equal(t, []string{"c", "d"}, values)
	})
}
