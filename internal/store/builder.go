package store

type Builder struct {
	clock Clock
}

func NewBuilder() Builder {
	return Builder{clock: SystemClock{}}
}

func (b Builder) WithClock(c Clock) Builder {
	b.clock = c
	return b
}

func (b Builder) Build() (Store, *ExpiryScanner) {
	tracker := NewExpiryTracker()
	dataStore := New().WithClock(b.clock).WithExpiryTracker(tracker)

	scanner := NewExpiryScanner(tracker, dataStore)

	return dataStore, scanner
}
