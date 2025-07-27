package list_test

import (
	"github.com/stretchr/testify/assert"
	"redis-challenge/internal/store/list"
	"testing"
)

func TestReadRangeFromStoreList(t *testing.T) {

	t.Run("range from empty list is an empty list", func(t *testing.T) {
		values, ok := list.ReadRangeFromStoreList(nil, 0, 10)

		assert.True(t, ok, "should be ok")
		assert.Empty(t, values, "should return an empty list")
	})

	t.Run("range from string is not ok", func(t *testing.T) {
		_, ok := list.ReadRangeFromStoreList("", 0, 10)

		assert.False(t, ok)
	})

	t.Run("range with start in range equal to end returns value", func(t *testing.T) {
		values, ok := list.ReadRangeFromStoreList([]string{"a", "b", "c", "d"}, 1, 1)

		assert.True(t, ok, "should be ok")
		assert.Equal(t, []string{"b"}, values)
	})

	t.Run("range with start to end both in range", func(t *testing.T) {
		values, ok := list.ReadRangeFromStoreList([]string{"a", "b", "c", "d", "e"}, 1, 3)

		assert.True(t, ok, "should be ok")
		assert.Equal(t, []string{"b", "c", "d"}, values)
	})

	t.Run("range with end after start both in range returns nil", func(t *testing.T) {
		values, ok := list.ReadRangeFromStoreList([]string{"a", "b", "c", "d", "e"}, 3, 1)

		assert.True(t, ok, "should be ok")
		assert.Nil(t, values)
	})

	t.Run("range with start in range to end beyond range", func(t *testing.T) {
		values, ok := list.ReadRangeFromStoreList([]string{"a", "b", "c", "d", "e"}, 3, 10)

		assert.True(t, ok, "should be ok")
		assert.Equal(t, []string{"d", "e"}, values, "read values up to the end of the list")
	})

	t.Run("range with start negative and end positive such that start is before end", func(t *testing.T) {
		values, ok := list.ReadRangeFromStoreList([]string{"a", "b", "c", "d", "e"}, -3, 3)

		assert.True(t, ok, "should be ok")
		assert.Equal(t, []string{"c", "d"}, values)
	})

	t.Run("range with start negative and end positive such that start counts back beyond list start is treated as back to list start", func(t *testing.T) {
		values, ok := list.ReadRangeFromStoreList([]string{"a", "b", "c", "d", "e"}, -100, 3)

		assert.True(t, ok, "should be ok")
		assert.Equal(t, []string{"a", "b", "c", "d"}, values)
	})

	t.Run("range with end -1 includes the end of the list", func(t *testing.T) {
		values, ok := list.ReadRangeFromStoreList([]string{"a", "b", "c", "d", "e"}, 1, -1)

		assert.True(t, ok, "should be ok")
		assert.Equal(t, []string{"b", "c", "d", "e"}, values)
	})

	t.Run("range with start positive and end negative such that start is before end", func(t *testing.T) {
		values, ok := list.ReadRangeFromStoreList([]string{"a", "b", "c", "d", "e"}, 1, -2)

		assert.True(t, ok, "should be ok")
		assert.Equal(t, []string{"b", "c", "d"}, values)
	})

	t.Run("range with start positive and end negative such that end is before start", func(t *testing.T) {
		values, ok := list.ReadRangeFromStoreList([]string{"a", "b", "c", "d", "e"}, 4, -3)

		assert.True(t, ok, "should be ok")
		assert.Empty(t, values)
	})
}
