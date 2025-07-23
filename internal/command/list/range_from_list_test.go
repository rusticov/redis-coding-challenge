package list_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"redis-challenge/internal/command/list"
	"testing"
)

func TestReadRangeFromStoreList(t *testing.T) {

	t.Run("range from empty list is an empty list", func(t *testing.T) {
		values, err := list.ReadRangeFromStoreList(nil, 0, 10)

		require.Nil(t, err, "should not return an error")
		assert.Empty(t, values, "should return an empty list")
	})

	t.Run("range with start in range equal to end returns value", func(t *testing.T) {
		values, err := list.ReadRangeFromStoreList([]string{"a", "b", "c", "d"}, 1, 1)

		require.Nil(t, err, "should not return an error")
		assert.Equal(t, []string{"b"}, values)
	})

	t.Run("range with start to end both in range", func(t *testing.T) {
		values, err := list.ReadRangeFromStoreList([]string{"a", "b", "c", "d", "e"}, 1, 3)

		require.Nil(t, err, "should not return an error")
		assert.Equal(t, []string{"b", "c", "d"}, values)
	})

	t.Run("range with end after start both in range returns nil", func(t *testing.T) {
		values, err := list.ReadRangeFromStoreList([]string{"a", "b", "c", "d", "e"}, 3, 1)

		require.Nil(t, err, "should not return an error")
		assert.Nil(t, values)
	})

	t.Run("range with start in range to end beyond range", func(t *testing.T) {
		values, err := list.ReadRangeFromStoreList([]string{"a", "b", "c", "d", "e"}, 3, 10)

		require.Nil(t, err, "should not return an error")
		assert.Equal(t, []string{"d", "e"}, values, "read values up to the end of the list")
	})

	t.Run("range with start negative and end positive such that start is before end", func(t *testing.T) {
		values, err := list.ReadRangeFromStoreList([]string{"a", "b", "c", "d", "e"}, -3, 3)

		require.Nil(t, err, "should not return an error")
		assert.Equal(t, []string{"c", "d"}, values)
	})

	t.Run("range with start negative and end positive such that start counts back beyond list start is treated as back to list start", func(t *testing.T) {
		values, err := list.ReadRangeFromStoreList([]string{"a", "b", "c", "d", "e"}, -100, 3)

		require.Nil(t, err, "should not return an error")
		assert.Equal(t, []string{"a", "b", "c", "d"}, values)
	})

	t.Run("range with start positive and end negative such that start is before end", func(t *testing.T) {
		values, err := list.ReadRangeFromStoreList([]string{"a", "b", "c", "d", "e"}, 1, -2)

		require.Nil(t, err, "should not return an error")
		assert.Equal(t, []string{"b", "c"}, values)
	})

	t.Run("range with start positive and end negative such that end is before start", func(t *testing.T) {
		values, err := list.ReadRangeFromStoreList([]string{"a", "b", "c", "d", "e"}, 4, -3)

		require.Nil(t, err, "should not return an error")
		assert.Empty(t, values)
	})
}
