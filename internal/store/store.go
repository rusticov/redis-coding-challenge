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
	values    *sync.Map
	keyValues map[string]Entry
}

func (s *InMemoryStore) exists(key string) bool {
	_, ok := s.keyValues[key]
	return ok
}

func (s *InMemoryStore) CompareAndSwap(key string, oldValue, newValue Entry) (swapped bool) {
	swapped = s.values.CompareAndSwap(key, oldValue, newValue)
	if s.keyValues[key] == oldValue {
		s.keyValues[key] = newValue
	}
	return
}

func (s *InMemoryStore) LoadOrStore(key string, defaultValue Entry) (Entry, bool) {
	value, loaded := s.values.LoadOrStore(key, defaultValue)
	if !s.exists(key) {
		s.keyValues[key] = defaultValue
	}
	return value.(Entry), loaded
}

func (s *InMemoryStore) Delete(key string) bool {
	s.values.LoadAndDelete(key)
	existed := s.exists(key)
	delete(s.keyValues, key)
	return existed
}

func (s *InMemoryStore) Get(key string) (Entry, bool) {
	value, ok := s.keyValues[key]
	return value, ok
}

func New() *InMemoryStore {
	return &InMemoryStore{
		values:    &sync.Map{},
		keyValues: make(map[string]Entry),
	}
}
