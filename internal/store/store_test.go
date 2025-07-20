package store_test

import (
	"github.com/stretchr/testify/assert"
	"redis-challenge/internal/store"
	"testing"
)

func TestStore(t *testing.T) {

	t.Run("set values can be read back", func(t *testing.T) {
		s := store.New()

		s.Add("key", "value")
		value, exists := s.Get("key")

		assert.Equal(t, "value", value)
		assert.True(t, exists)
	})

	t.Run("get value for non-existent key declares key does not have a value", func(t *testing.T) {
		s := store.New()

		_, exists := s.Get("key")

		assert.False(t, exists)
	})
}
