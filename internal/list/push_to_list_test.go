package list_test

import (
	"github.com/stretchr/testify/assert"
	"redis-challenge/internal/list"
	"testing"
)

func TestPushToList(t *testing.T) {

	t.Run("twice left push string to non-empty list of strings maintains order", func(t *testing.T) {
		pushedList, ok := list.LeftPush([]string{"a", "b", "c"}, nil)
		assert.True(t, ok, "should return ok")

		pushedList, ok = list.LeftPush([]string{"1", "2", "3"}, pushedList)

		assert.True(t, ok, "should return ok")

		values, ok := list.ReadRangeFromStoreList(pushedList, 0, -1)
		assert.Equal(t, []string{"3", "2", "1", "c", "b", "a"}, values, "only new values are reversed")
	})

	t.Run("twice right push string to non-empty list of strings maintains order", func(t *testing.T) {
		pushedList, ok := list.RightPush([]string{"d", "e", "f"}, nil)
		assert.True(t, ok, "should be ok")

		pushedList, ok = list.RightPush([]string{"a", "b", "c"}, pushedList)

		assert.True(t, ok, "should be ok")

		values, ok := list.ReadRangeFromStoreList(pushedList, 0, -1)
		assert.Equal(t, []string{"d", "e", "f", "a", "b", "c"}, values, "value are in order")
	})
}
