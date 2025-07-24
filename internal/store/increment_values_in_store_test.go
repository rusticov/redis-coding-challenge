package store_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"redis-challenge/internal/store"
	"testing"
)

func TestIncrementValuesInStore(t *testing.T) {

	t.Run("increment a number of an unstored key", func(t *testing.T) {
		s := store.New()

		updatedValue, err := s.Increment("key", int64(3))
		require.Nil(t, err)

		assert.Equal(t, int64(3), updatedValue)
	})

	t.Run("increment a number of key with integer", func(t *testing.T) {
		s := store.New()
		s.Write("key", "10")

		updatedValue, err := s.Increment("key", int64(3))
		require.Nil(t, err)

		assert.Equal(t, int64(13), updatedValue)
	})

	t.Run("read the incremented number of an unstored key", func(t *testing.T) {
		s := store.New()
		s.Write("key", "42")

		_, err := s.Increment("key", int64(-2))
		require.Nil(t, err)

		value, err := s.ReadString("key")

		require.Nil(t, err)
		assert.Equal(t, "40", value)
	})

	t.Run("increment a number of key with a non-integer", func(t *testing.T) {
		s := store.New()
		s.Write("key", "ten")

		_, err := s.Increment("key", int64(3))

		assert.Equal(t, store.ErrorNotAnInteger, err)
	})
}
