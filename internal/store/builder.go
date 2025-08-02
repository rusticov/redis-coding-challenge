package store

import "io"

type Builder struct {
	clock            Clock
	commandLogWriter io.Writer
}

func NewBuilder() Builder {
	return Builder{clock: SystemClock{}}
}

func (b Builder) WithClock(c Clock) Builder {
	b.clock = c
	return b
}

func (b Builder) WithCommandLogWriter(writer io.Writer) Builder {
	b.commandLogWriter = writer
	return b
}

func (b Builder) Build() (Store, *ExpiryScanner) {
	tracker := NewExpiryTracker().withDeleteListener(&deleteListener{writer: b.commandLogWriter})
	dataStore := New().WithClock(b.clock).WithExpiryTracker(tracker)

	scanner := NewExpiryScanner(tracker, dataStore)

	return dataStore, scanner
}
