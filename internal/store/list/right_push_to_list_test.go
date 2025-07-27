package list_test

import (
	"github.com/stretchr/testify/assert"
	"redis-challenge/internal/store/list"
	"testing"
)

func TestRightPushValuesToList(t *testing.T) {

	t.Run("right push string to empty list", func(t *testing.T) {
		newList, ok := list.RightPushToOldList([]string{"a"}, nil)

		assert.True(t, ok, "should be ok")
		assert.Equal(t, []string{"a"}, newList)
	})

	t.Run("right push 2 strings to empty list should be in order", func(t *testing.T) {
		newList, ok := list.RightPushToOldList([]string{"a", "b"}, nil)

		assert.True(t, ok, "should be ok")
		assert.Equal(t, []string{"a", "b"}, newList)
	})

	t.Run("right push string to non-empty list of strings", func(t *testing.T) {
		newList, ok := list.RightPushToOldList([]string{"a", "b", "c"}, []string{"d", "e", "f"})

		assert.True(t, ok, "should be ok")
		assert.Equal(t, []string{"d", "e", "f", "a", "b", "c"}, newList, "value are in order")
	})

	t.Run("right push string to old value that is not a list should error", func(t *testing.T) {
		_, ok := list.RightPushToOldList([]string{"a", "b", "c"}, "f")

		assert.False(t, ok, "should not be ok")
	})
}
