package list_test

import (
	"github.com/stretchr/testify/assert"
	"redis-challenge/internal/list"
	"slices"
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

	testCases := map[string]struct {
		storedList []string
		start      int
		end        int
		expected   []string
	}{
		"range from empty list is an empty list": {
			storedList: nil,
			start:      0,
			end:        10,
			expected:   nil,
		},
		"range from right-pushed list with one value is the value": {
			storedList: []string{"a", "b", "c", "d"},
			start:      1,
			end:        1,
			expected:   []string{"b"},
		},
		"range from right-pushed list with start and end both in range": {
			storedList: []string{"a", "b", "c", "d", "e"},
			start:      1,
			end:        3,
			expected:   []string{"b", "c", "d"},
		},
		"range from right-pushed list with end after start both in range returns nil": {
			storedList: []string{"a", "b", "c", "d", "e"},
			start:      3,
			end:        1,
			expected:   nil,
		},
		"range from right-pushed list with start in range and end beyond range": {
			storedList: []string{"a", "b", "c", "d", "e"},
			start:      3,
			end:        10,
			expected:   []string{"d", "e"},
		},
		"range from right-pushed list with start negative and end positive such that start is before end": {
			storedList: []string{"a", "b", "c", "d", "e"},
			start:      -3,
			end:        3,
			expected:   []string{"c", "d"},
		},
		"range from right-pushed list with start negative and end positive such that start counts back beyond list start is treated as back to list start": {
			storedList: []string{"a", "b", "c", "d", "e"},
			start:      -100,
			end:        3,
			expected:   []string{"a", "b", "c", "d"},
		},
		"range from right-pushed list with end -1 includes the end of the list": {
			storedList: []string{"a", "b", "c", "d", "e"},
			start:      1,
			end:        -1,
			expected:   []string{"b", "c", "d", "e"},
		},
		"range from right-pushed list with start positive and end negative such that start is before end": {
			storedList: []string{"a", "b", "c", "d", "e"},
			start:      1,
			end:        -2,
			expected:   []string{"b", "c", "d"},
		},
		"range from right-pushed list with start positive and end negative such that end is before start": {
			storedList: []string{"a", "b", "c", "d", "e"},
			start:      4,
			end:        -3,
			expected:   nil,
		},
	}

	for name, testCase := range testCases {
		t.Run("range from right-pushed list "+name, func(t *testing.T) {
			rightPushedList, _ := list.RightPush(testCase.storedList, nil)

			values, ok := list.ReadRangeFromStoreList(rightPushedList, testCase.start, testCase.end)

			assert.True(t, ok, "should be ok")
			assert.Equal(t, testCase.expected, values)
		})

		t.Run("range from left-pushed list "+name, func(t *testing.T) {
			xs := make([]string, len(testCase.storedList))
			copy(xs, testCase.storedList)
			slices.Reverse(xs)

			leftPushedList, _ := list.LeftPush(xs, nil)

			values, ok := list.ReadRangeFromStoreList(leftPushedList, testCase.start, testCase.end)

			assert.True(t, ok, "should be ok")
			assert.Equal(t, testCase.expected, values)
		})

		t.Run("iterate from right-pushed list "+name, func(t *testing.T) {
			rightPushedList, _ := list.RightPush(testCase.storedList, nil)
			rangeLimitedList := rightPushedList.Filter(testCase.start, testCase.end)

			var result []string
			for _, value := range rangeLimitedList.Range() {
				result = append(result, value)
			}

			assert.Equal(t, testCase.expected, result)
		})
	}

	t.Run("range from left and right pushed list with all values found on the right", func(t *testing.T) {
		pushedList, _ := list.LeftPush([]string{"1", "2", "3"}, nil)
		pushedList, _ = list.RightPush([]string{"a", "b", "c", "d", "e"}, pushedList)

		values, ok := list.ReadRangeFromStoreList(pushedList, 4, 6)

		assert.True(t, ok, "should be ok")
		assert.Equal(t, []string{"b", "c", "d"}, values)
	})

	t.Run("range from left and right pushed list with values that straddle both lists", func(t *testing.T) {
		pushedList, _ := list.LeftPush([]string{"1", "2", "3"}, nil)
		pushedList, _ = list.RightPush([]string{"a", "b", "c", "d", "e"}, pushedList)

		values, ok := list.ReadRangeFromStoreList(pushedList, 1, 4)

		assert.True(t, ok, "should be ok")
		assert.Equal(t, []string{"2", "1", "a", "b"}, values)
	})
}
