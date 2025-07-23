package store

import (
	"sync"
)

type Store interface {
	Get(key string) (string, bool)
	CompareAndSwap(key string, oldValue, newValue string) (swapped bool)
	LoadOrStore(key string, defaultValue string) (string, bool)
	Delete(key string) bool
}

type InMemoryStore struct {
	values *sync.Map
}

func (s *InMemoryStore) CompareAndSwap(key string, oldValue string, newValue string) (swapped bool) {
	swapped = s.values.CompareAndSwap(key, oldValue, newValue)
	return
}

func (s *InMemoryStore) LoadOrStore(key string, defaultValue string) (string, bool) {
	value, loaded := s.values.LoadOrStore(key, defaultValue)
	return value.(string), loaded
}

func (s *InMemoryStore) Delete(key string) bool {
	_, existed := s.values.LoadAndDelete(key)
	return existed
}

func (s *InMemoryStore) Get(key string) (string, bool) {
	value, ok := s.values.Load(key)

	if !ok {
		return "", false
	}
	return value.(string), true
}

func New() *InMemoryStore {
	return &InMemoryStore{values: &sync.Map{}}
}
