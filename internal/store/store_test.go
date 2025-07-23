package store_test

import (
	"github.com/stretchr/testify/assert"
	"redis-challenge/internal/store"
	"testing"
)

func TestStore(t *testing.T) {

	t.Run("newly sets values can be read back", func(t *testing.T) {
		s := store.New()

		value, loaded := s.LoadOrStore("key", store.NewEntry("value"))
		assert.False(t, loaded)
		assert.Equal(t, store.NewEntry("value"), value)

		value, exists := s.Get("key")

		assert.Equal(t, store.NewEntry("value"), value)
		assert.True(t, exists)
	})

	t.Run("setting a value that is already set will load existing value", func(t *testing.T) {
		s := store.New()

		s.LoadOrStore("key", store.NewEntry("value 1"))
		value, loaded := s.LoadOrStore("key", store.NewEntry("value 2"))
		assert.True(t, loaded)
		assert.Equal(t, store.NewEntry("value 1"), value, "should return the existing value as one has been set")

		value, exists := s.Get("key")

		assert.Equal(t, store.NewEntry("value 1"), value)
		assert.True(t, exists)
	})

	t.Run("values can be updated if the current value is given as the old value", func(t *testing.T) {
		s := store.New()

		s.LoadOrStore("key", store.NewEntry("value 1"))
		s.CompareAndSwap("key", store.NewEntry("value 1"), store.NewEntry("value 2"))
		value, exists := s.Get("key")

		assert.Equal(t, store.NewEntry("value 2"), value)
		assert.True(t, exists)
	})

	t.Run("values is not updated if the old value is different to the current value", func(t *testing.T) {
		s := store.New()

		s.LoadOrStore("key", store.NewEntry("value 1"))
		s.CompareAndSwap("key", store.NewEntry("not current value"), store.NewEntry("value 1"))
		value, exists := s.Get("key")

		assert.Equal(t, store.NewEntry("value 1"), value)
		assert.True(t, exists)
	})

	t.Run("get value for non-existent key declares key does not have a value", func(t *testing.T) {
		s := store.New()

		_, exists := s.Get("key")

		assert.False(t, exists)
	})
}
