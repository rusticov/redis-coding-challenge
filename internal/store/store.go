package store

import "strconv"

type Store interface {
	Get(key string) (Entry, bool)
	CompareAndSwap(key string, oldValue, newValue Entry) (swapped bool)
	LoadOrStore(key string, defaultValue Entry) (Entry, bool)
	Delete(key string) bool
}

type InMemoryStore struct {
	keyValues map[string]Entry
}

func (s *InMemoryStore) exists(key string) bool {
	_, ok := s.keyValues[key]
	return ok
}

func (s *InMemoryStore) CompareAndSwap(key string, oldValue, newValue Entry) (swapped bool) {
	swapped = s.keyValues[key] == oldValue
	if swapped {
		s.keyValues[key] = newValue
	}
	return
}

func (s *InMemoryStore) LoadOrStore(key string, defaultValue Entry) (Entry, bool) {
	loaded := s.exists(key)
	if !loaded {
		s.keyValues[key] = defaultValue
	}
	return s.keyValues[key], loaded
}

func (s *InMemoryStore) Delete(key string) bool {
	existed := s.exists(key)
	delete(s.keyValues, key)
	return existed
}

func (s *InMemoryStore) Get(key string) (Entry, bool) {
	value, ok := s.keyValues[key]
	return value, ok
}

func (s *InMemoryStore) Write(key string, value any) {
	s.keyValues[key] = NewEntry(value)
}

func (s *InMemoryStore) ReadString(key string) (string, error) {
	if value, ok := s.keyValues[key]; ok {
		return value.data.(string), nil
	}
	return "", ErrorKeyNotFound
}

func (s *InMemoryStore) Increment(key string, incrementBy int64) (int64, error) {
	var value int64
	if entry, ok := s.keyValues[key]; ok {
		var err error
		value, err = strconv.ParseInt(entry.data.(string), 10, 64)
		if err != nil {
			return 0, ErrorNotAnInteger
		}
	}
	value += incrementBy

	stringValue := strconv.FormatInt(value, 10)
	s.keyValues[key] = NewEntry(stringValue)
	return value, nil
}

func New() *InMemoryStore {
	return &InMemoryStore{
		keyValues: make(map[string]Entry),
	}
}
