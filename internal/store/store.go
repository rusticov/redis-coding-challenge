package store

import (
	"sync"
)

type Store interface {
	Get(key string) (Entry, bool)
	CompareAndSwap(key string, oldValue, newValue Entry) (swapped bool)
	LoadOrStore(key string, defaultValue Entry) (Entry, bool)
	Delete(key string) bool
}

type InMemoryStore struct {
	values *sync.Map
}

func (s *InMemoryStore) CompareAndSwap(key string, oldValue, newValue Entry) (swapped bool) {
	swapped = s.values.CompareAndSwap(key, oldValue, newValue)
	return
}

func (s *InMemoryStore) LoadOrStore(key string, defaultValue Entry) (Entry, bool) {
	value, loaded := s.values.LoadOrStore(key, defaultValue)
	return value.(Entry), loaded
}

func (s *InMemoryStore) Delete(key string) bool {
	_, existed := s.values.LoadAndDelete(key)
	return existed
}

func (s *InMemoryStore) Get(key string) (Entry, bool) {
	value, ok := s.values.Load(key)

	if !ok {
		return Entry{}, false
	}
	return value.(Entry), true
}

func New() *InMemoryStore {
	return &InMemoryStore{values: &sync.Map{}}
}
