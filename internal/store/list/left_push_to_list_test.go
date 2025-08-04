package list_test

import (
	"github.com/stretchr/testify/assert"
	"redis-challenge/internal/store/list"
	"testing"
)

func TestAddValuesToList(t *testing.T) {

	t.Run("left push string to empty list", func(t *testing.T) {
		newList, ok := list.LeftPushToOldList([]string{"a"}, nil)

		assert.True(t, ok, "should return ok")
		assert.Equal(t, []string{"a"}, newList.Left)
	})

	t.Run("left push 2 strings to empty list should reverse order", func(t *testing.T) {
		newList, ok := list.LeftPushToOldList([]string{"a", "b"}, nil)

		assert.True(t, ok, "should return ok")
		assert.Equal(t, []string{"a", "b"}, newList.Left)
	})

	t.Run("left push string to non-empty list of strings", func(t *testing.T) {
		newList, ok := list.LeftPushToOldList([]string{"a", "b", "c"}, list.DoubleEndedList{Left: []string{"d", "e", "f"}})

		assert.True(t, ok, "should return ok")
		assert.Equal(t, []string{"d", "e", "f", "a", "b", "c"}, newList.Left, "only new values are reversed")
	})

	t.Run("left push string to old value that is not a list should error", func(t *testing.T) {
		_, ok := list.LeftPushToOldList([]string{"a", "b", "c"}, "f")

		assert.False(t, ok, "should not be ok")
	})
}
