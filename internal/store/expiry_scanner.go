package store

type ExpiryScanner struct {
	tracker            *ExpiryTracker
	store              Store
	randomCount        int
	continuePurgeCount int
}

func (s ExpiryScanner) Scan() {
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

func NewExpiryScanner(tracker *ExpiryTracker, s Store) *ExpiryScanner {
	return &ExpiryScanner{
		tracker:            tracker,
		store:              s,
		randomCount:        20,
		continuePurgeCount: 6, // if more than 25% are purged
	}
}
