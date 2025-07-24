package store

import "strconv"

type Store interface {
	Get(key string) (Entry, bool)
	CompareAndSwap(key string, oldValue, newValue Entry) (swapped bool)
	LoadOrStore(key string, defaultValue Entry) (Entry, bool)
	Delete(key string) bool

	Write(key string, value any)
	ReadString(key string) (string, error)
	Exists(key string) bool
	Increment(key string, incrementBy int64) (int64, error)
}

type InMemoryStore struct {
	keyValues map[string]Entry
}

func (s *InMemoryStore) Exists(key string) bool {
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
	loaded := s.Exists(key)
	if !loaded {
		s.keyValues[key] = defaultValue
	}
	return s.keyValues[key], loaded
}

func (s *InMemoryStore) Delete(key string) bool {
	existed := s.Exists(key)
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
	value, err := s.readInteger(key)
	if err != nil {
		return 0, err
	}
	value += incrementBy

	stringValue := strconv.FormatInt(value, 10)
	s.keyValues[key] = NewEntry(stringValue)
	return value, nil
}

func (s *InMemoryStore) readInteger(key string) (int64, error) {
	if entry, ok := s.keyValues[key]; ok {
		value, err := strconv.ParseInt(entry.data.(string), 10, 64)
		if err != nil {
			return 0, ErrorNotAnInteger
		}
		return value, nil
	}
	return 0, nil
}

func New() *InMemoryStore {
	return &InMemoryStore{
		keyValues: make(map[string]Entry),
	}
}
