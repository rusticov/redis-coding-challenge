package store

import (
	"sync"
)

type Store struct {
	values *sync.Map
}

func (s *Store) Add(key string, value string) {
	s.values.Store(key, value)
}

func (s *Store) Get(key string) (string, bool) {
	value, ok := s.values.Load(key)

	if !ok {
		return "", false
	}
	return value.(string), true
}

func New() *Store {
	return &Store{values: &sync.Map{}}
}
