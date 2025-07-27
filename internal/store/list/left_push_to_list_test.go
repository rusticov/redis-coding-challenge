package list_test

import (
	"github.com/stretchr/testify/assert"
	"redis-challenge/internal/store/list"
	"testing"
)

func TestAddValuesToList(t *testing.T) {

	t.Run("left push string to empty list", func(t *testing.T) {
		newList, errorData := list.LeftPushToOldList([]string{"a"}, nil)

		assert.Nil(t, errorData, "should not return an error")
		assert.Equal(t, []string{"a"}, newList)
	})

	t.Run("left push 2 strings to empty list should reverse order", func(t *testing.T) {
		newList, errorData := list.LeftPushToOldList([]string{"a", "b"}, nil)

		assert.Nil(t, errorData, "should not return an error")
		assert.Equal(t, []string{"b", "a"}, newList)
	})

	t.Run("left push string to non-empty list of strings", func(t *testing.T) {
		newList, errorData := list.LeftPushToOldList([]string{"a", "b", "c"}, []string{"d", "e", "f"})

		assert.Nil(t, errorData, "should not return an error")
		assert.Equal(t, []string{"c", "b", "a", "d", "e", "f"}, newList, "only new values are reversed")
	})

	t.Run("left push string to old value that is not a list should error", func(t *testing.T) {
		_, errorData := list.LeftPushToOldList([]string{"a", "b", "c"}, "f")

		assert.Equal(t, list.ErrorOldValueIsNotList, errorData, "should return an error")
	})
}
