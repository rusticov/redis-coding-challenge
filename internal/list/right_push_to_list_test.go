package list_test

import (
	"github.com/stretchr/testify/assert"
	list2 "redis-challenge/internal/list"
	"testing"
)

func TestRightPushValuesToList(t *testing.T) {

	t.Run("right push string to empty list", func(t *testing.T) {
		newList, ok := list2.RightPush([]string{"a"}, nil)

		assert.True(t, ok, "should be ok")
		assert.Equal(t, []string{"a"}, newList.right)
	})

	t.Run("right push 2 strings to empty list should be in order", func(t *testing.T) {
		newList, ok := list2.RightPush([]string{"a", "b"}, nil)

		assert.True(t, ok, "should be ok")
		assert.Equal(t, []string{"a", "b"}, newList.right)
	})

	t.Run("right push string to non-empty list of strings", func(t *testing.T) {
		newList, ok := list2.RightPush([]string{"a", "b", "c"}, list2.DoubleEndedList{right: []string{"d", "e", "f"}})

		assert.True(t, ok, "should be ok")
		assert.Equal(t, []string{"d", "e", "f", "a", "b", "c"}, newList.right, "value are in order")
	})

	t.Run("right push string to old value that is not a list should error", func(t *testing.T) {
		_, ok := list2.RightPush([]string{"a", "b", "c"}, "f")

		assert.False(t, ok, "should not be ok")
	})
}
