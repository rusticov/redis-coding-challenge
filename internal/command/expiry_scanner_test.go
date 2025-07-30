package command_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"redis-challenge/internal/command"
	"redis-challenge/internal/store"
	"testing"
)

func TestExpiryScanner(t *testing.T) {

	t.Run("purge key when it has become expired", func(t *testing.T) {
		clock := store.FixedClock{}
		tracker := store.NewExpiryTracker()

		s := store.NewWithClock(clock.Now).WithExpiryTracker(tracker)

		scanner := command.NewExpiryScanner(tracker, s)

		s.Write("key", "value", store.ExpiryOptionExpirySeconds, 1)

		clock.AddSeconds(2)

		scanner.Scan()

		assert.Empty(t, tracker.SelectKeys(1), "no more keys marked as expired")
		assert.Equal(t, 0, s.Size(), "no more keys in store after only key is purged")
	})

	t.Run("not successfully purge key when it has an expiry but is it not expired", func(t *testing.T) {
		clock := store.FixedClock{}
		tracker := store.NewExpiryTracker()

		s := store.NewWithClock(clock.Now).WithExpiryTracker(tracker)

		scanner := command.NewExpiryScanner(tracker, s)

		s.Write("key", "value", store.ExpiryOptionExpirySeconds, 1)

		scanner.Scan()

		assert.Contains(t, tracker.SelectKeys(1), "key", "key is still tracked as it has an expiry")
		assert.True(t, s.Exists("key"), "key is still in store as it has not yet expired")
	})

	t.Run("purge all expired keys if random selection includes all keys with expiry", func(t *testing.T) {
		clock := store.FixedClock{}
		tracker := store.NewExpiryTracker()

		s := store.NewWithClock(clock.Now).WithExpiryTracker(tracker)

		scanner := command.NewExpiryScanner(tracker, s)

		s.Write("key1", "value", store.ExpiryOptionExpirySeconds, 1)
		s.Write("key2", "value", store.ExpiryOptionExpirySeconds, 4)
		s.Write("key3", "value", store.ExpiryOptionExpirySeconds, 2)
		s.Write("key4", "value", store.ExpiryOptionExpirySeconds, 5)

		clock.AddSeconds(3)

		scanner.Scan()

		remainingSelection := tracker.SelectKeys(10)

		assert.Len(t, remainingSelection, 2, "2 keys remain as they have not yet expired")
		assert.Contains(t, remainingSelection, "key2", "key with expiry in the future is still in the selection")
		assert.Contains(t, remainingSelection, "key4", "key with expiry in the future is still in the selection")

		assert.Equal(t, 2, s.Size(), "expired keys are removed from store")
	})

	t.Run("purge more than count when a purge fraction exceeds defined limit", func(t *testing.T) {
		clock := store.FixedClock{}
		tracker := store.NewExpiryTracker()

		s := store.NewWithClock(clock.Now).WithExpiryTracker(tracker)

		scanner := command.NewExpiryScanner(tracker, s)

		for i := range 100 {
			s.Write(fmt.Sprintf("key%d", i), "value", store.ExpiryOptionExpirySeconds, 1)
		}
		s.Write("unexpired-key", "value", store.ExpiryOptionExpirySeconds, 15)

		clock.AddSeconds(10)

		scanner.Scan()

		remainingSelection := tracker.SelectKeys(10)

		assert.Less(t, len(remainingSelection), 10, "most keys purged")
		assert.Contains(t, remainingSelection, "unexpired-key", "key with expiry in the future is still in the selection")

		assert.Equal(t, len(remainingSelection), s.Size(), "expired keys are removed from store")
	})
}
