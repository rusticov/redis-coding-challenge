package list_test

import (
	"github.com/stretchr/testify/require"
	list2 "redis-challenge/internal/store/list"
	"testing"
)

func TestRightPushValuesToList(t *testing.T) {

	t.Run("right push string to empty list", func(t *testing.T) {
		newList, errorData := list2.RightPushToOldList([]string{"a"}, nil)

		require.Nil(t, errorData, "should not return an error")
		require.Equal(t, []string{"a"}, newList)
	})

	t.Run("right push 2 strings to empty list should be in order", func(t *testing.T) {
		newList, errorData := list2.RightPushToOldList([]string{"a", "b"}, nil)

		require.Nil(t, errorData, "should not return an error")
		require.Equal(t, []string{"a", "b"}, newList)
	})

	t.Run("right push string to non-empty list of strings", func(t *testing.T) {
		newList, errorData := list2.RightPushToOldList([]string{"a", "b", "c"}, []string{"d", "e", "f"})

		require.Nil(t, errorData, "should not return an error")
		require.Equal(t, []string{"d", "e", "f", "a", "b", "c"}, newList, "value are in order")
	})

	t.Run("right push string to old value that is not a list should error", func(t *testing.T) {
		_, errorData := list2.RightPushToOldList([]string{"a", "b", "c"}, "f")

		require.Equal(t, list2.ErrorOldValueIsNotList, errorData, "should return an error")
	})
}
