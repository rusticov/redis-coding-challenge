package list

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddValuesToList(t *testing.T) {

	t.Run("left push string to empty list", func(t *testing.T) {
		newList, ok := LeftPush([]string{"a"}, nil)

		assert.True(t, ok, "should return ok")
		assert.Equal(t, []string{"a"}, newList.left)
	})

	t.Run("left push 2 strings to empty list should reverse order", func(t *testing.T) {
		newList, ok := LeftPush([]string{"a", "b"}, nil)

		assert.True(t, ok, "should return ok")
		assert.Equal(t, []string{"a", "b"}, newList.left)
	})

	t.Run("left push string to non-empty list of strings", func(t *testing.T) {
		newList, ok := LeftPush([]string{"a", "b", "c"}, DoubleEndedList{left: []string{"d", "e", "f"}})

		assert.True(t, ok, "should return ok")
		assert.Equal(t, []string{"d", "e", "f", "a", "b", "c"}, newList.left, "only new values are reversed")
	})

	t.Run("left push string to old value that is not a list should error", func(t *testing.T) {
		_, ok := LeftPush([]string{"a", "b", "c"}, "f")

		assert.False(t, ok, "should not be ok")
	})
}
