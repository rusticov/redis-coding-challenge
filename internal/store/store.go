package store

import (
	"sync"
)

type Store interface {
	Get(key string) (any, bool)
	CompareAndSwap(key string, oldValue, newValue any) (swapped bool)
	LoadOrStore(key string, defaultValue any) (any, bool)
	Delete(key string) bool
}

type InMemoryStore struct {
	values *sync.Map
}

func (s *InMemoryStore) CompareAndSwap(key string, oldValue, newValue any) (swapped bool) {
	swapped = s.values.CompareAndSwap(key, oldValue, newValue)
	return
}

func (s *InMemoryStore) LoadOrStore(key string, defaultValue any) (any, bool) {
	value, loaded := s.values.LoadOrStore(key, defaultValue)
	return value, loaded
}

func (s *InMemoryStore) Delete(key string) bool {
	_, existed := s.values.LoadAndDelete(key)
	return existed
}

func (s *InMemoryStore) Get(key string) (any, bool) {
	value, ok := s.values.Load(key)

	if !ok {
		return nil, false
	}
	return value, true
}

func New() *InMemoryStore {
	return &InMemoryStore{values: &sync.Map{}}
}
