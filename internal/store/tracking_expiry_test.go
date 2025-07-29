package store_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"redis-challenge/internal/store"
	"testing"
)

func TestTrackingExpiry(t *testing.T) {

	t.Run("writing a key without expiry is not tracked", func(t *testing.T) {
		tracker := store.NewExpiryTracker()
		s := store.New().WithExpiryTracker(tracker)

		s.Write("key", "value", store.ExpiryOptionNone, 0)

		assert.Empty(t, tracker.SelectKeys(1), "should not have any keys")
	})

	t.Run("writing a key with expiry is tracked", func(t *testing.T) {
		tracker := store.NewExpiryTracker()
		s := store.New().WithExpiryTracker(tracker)

		s.Write("key", "value", store.ExpiryOptionExpirySeconds, 1)

		assert.Equal(t, []string{"key"}, tracker.SelectKeys(10), "should have key as it is only one with expiry")
	})

	t.Run("writing two keys with expiry are both tracked", func(t *testing.T) {
		tracker := store.NewExpiryTracker()
		s := store.New().WithExpiryTracker(tracker)

		s.Write("key1", "value", store.ExpiryOptionExpirySeconds, 1)
		s.Write("key2", "value", store.ExpiryOptionExpirySeconds, 1)

		keys := tracker.SelectKeys(10)
		assert.Len(t, keys, 2)
		assert.Contains(t, keys, "key1")
		assert.Contains(t, keys, "key2")
	})

	t.Run("updating a key with a different expiry is tracked once", func(t *testing.T) {
		tracker := store.NewExpiryTracker()
		s := store.New().WithExpiryTracker(tracker)

		s.Write("key", "value 1", store.ExpiryOptionExpirySeconds, 1)
		s.Write("key", "value 2", store.ExpiryOptionExpirySeconds, 2)

		assert.Equal(t, []string{"key"}, tracker.SelectKeys(10), "should have key as it is only one with expiry")
	})

	t.Run("updating a key from having expiry to not having expiry is not tracked", func(t *testing.T) {
		tracker := store.NewExpiryTracker()
		s := store.New().WithExpiryTracker(tracker)

		s.Write("key", "value 1", store.ExpiryOptionExpirySeconds, 1)
		s.Write("key", "value 2", store.ExpiryOptionNone, 0)

		assert.Empty(t, tracker.SelectKeys(10), "should have key removed as it is no longer with expiry")
	})

	t.Run("deleting a key with expiry is not tracked", func(t *testing.T) {
		tracker := store.NewExpiryTracker()
		s := store.New().WithExpiryTracker(tracker)

		s.Write("key", "value 1", store.ExpiryOptionExpirySeconds, 1)
		s.Delete("key")

		assert.Empty(t, tracker.SelectKeys(10), "should have key removed")
	})

	t.Run("random keys are selected up to requested count", func(t *testing.T) {
		tracker := store.NewExpiryTracker()
		s := store.New().WithExpiryTracker(tracker)

		for i := range 30 {
			s.Write(fmt.Sprintf("key%d", i), "value", store.ExpiryOptionExpirySeconds, 1)
		}

		selection1 := tracker.SelectKeys(10)
		selection2 := tracker.SelectKeys(10)

		assert.Len(t, selection1, 10, "should have requested number of keys if there are more keys")
		assert.Len(t, selection2, 10, "should have requested number of keys if there are more keys")

		assert.NotEqual(t, selection2, selection1, "selection should be different each time")
	})

	t.Run("keys within the requested count is unchanged by further tracking changes", func(t *testing.T) {
		tracker := store.NewExpiryTracker()
		s := store.New().WithExpiryTracker(tracker)

		for i := range 7 {
			s.Write(fmt.Sprintf("key%d", i), "value", store.ExpiryOptionExpirySeconds, 1)
		}

		selection := tracker.SelectKeys(10)

		s.Write("extra-key", "value", store.ExpiryOptionExpirySeconds, 1)

		assert.Len(t, selection, 7, "should not have extra key in selection")
	})
}
