package store

import "time"

type Clock interface {
	Now() int64
}

type SystemClock struct{}

func (c SystemClock) Now() int64 {
	return time.Now().UTC().UnixMilli()
}

type FixedClock struct {
	TimeInMilliseconds int64
}

func (c *FixedClock) Now() int64 {
	return c.TimeInMilliseconds
}

func (c *FixedClock) AddSeconds(delta int64) *FixedClock {
	c.TimeInMilliseconds += delta * 1000
	return c
}

func (c *FixedClock) AddMilliseconds(delta int64) *FixedClock {
	c.TimeInMilliseconds += delta
	return c
}
