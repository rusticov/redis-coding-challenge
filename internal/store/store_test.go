package store_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"redis-challenge/internal/store"
	"testing"
)

func TestStore(t *testing.T) {

	t.Run("read value added to store for a key", func(t *testing.T) {
		s := store.New()
		s.Write("key", "value")

		value, err := s.ReadString("key")

		require.Nil(t, err)
		assert.Equal(t, "value", value)
	})

	t.Run("read different value adding to store for a key", func(t *testing.T) {
		s := store.New()
		s.Write("key", "different value")

		value, err := s.ReadString("key")

		require.Nil(t, err)
		assert.Equal(t, "different value", value)
	})

	t.Run("read value of key not added to store", func(t *testing.T) {
		s := store.New()

		_, err := s.ReadString("key")

		assert.Equal(t, store.ErrorKeyNotFound, err)
	})

	t.Run("read value of key that is deleted cannot be found", func(t *testing.T) {
		s := store.New()
		s.Write("key", "value")
		s.Delete("key")

		_, err := s.ReadString("key")

		assert.Equal(t, store.ErrorKeyNotFound, err)
	})

	t.Run("writing value to a key that is already in the store overwrites the value", func(t *testing.T) {
		s := store.New()
		s.Write("key", "value")
		s.Write("key", "different value")
		value, err := s.ReadString("key")

		require.Nil(t, err)
		assert.Equal(t, "different value", value)
	})

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
