package command

import "redis-challenge/internal/store"

type ExpiryScanner struct {
	tracker            *store.ExpiryTracker
	store              store.Store
	randomCount        int
	continuePurgeCount int
}

func (s ExpiryScanner) PurgeExpiredKeys() {
	for {
		selection := s.tracker.SelectKeys(s.randomCount)

		count := 0
		for _, key := range selection {
			if !s.store.Exists(key) {
				count++
			}
		}

		if count < s.continuePurgeCount {
			return
		}
	}
}

func (s ExpiryScanner) WithRandomCount(count int) ExpiryScanner {
	s.randomCount = count
	return s
}

func (s ExpiryScanner) WithContinuePurgeCount(continuePurgeCount int) ExpiryScanner {
	s.continuePurgeCount = continuePurgeCount
	return s
}

func NewExpiryScanner(tracker *store.ExpiryTracker, s store.Store) *ExpiryScanner {
	return &ExpiryScanner{
		tracker:            tracker,
		store:              s,
		randomCount:        20,
		continuePurgeCount: 4,
	}
}
