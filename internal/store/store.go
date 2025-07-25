package store

import (
	"redis-challenge/internal/store/list"
	"strconv"
)

type Store interface {
	Get(key string) (Entry, bool)
	LoadOrStore(key string, defaultValue Entry) (Entry, bool)
	Delete(key string) bool

	ReadString(key string) (string, error)
	ReadListRange(key string, fromIndex int, toIndex int) ([]string, error)
	Exists(key string) bool

	Write(key string, value any)
	Increment(key string, incrementBy int64) (int64, error)
	LeftPush(key string, values []string) (int64, error)
	RightPush(key string, values []string) (int64, error)
}

type InMemoryStore struct {
	keyValues map[string]Entry
}

func (s *InMemoryStore) Exists(key string) bool {
	_, ok := s.keyValues[key]
	return ok
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

func (s *InMemoryStore) LeftPush(key string, values []string) (int64, error) {
	oldList, _ := s.keyValues[key]
	updatedList, err := list.LeftPushToOldList(values, oldList.data)
	if err != nil {
		return 0, err
	}

	s.keyValues[key] = NewEntry(updatedList)

	return int64(len(updatedList)), nil
}

func (s *InMemoryStore) RightPush(key string, values []string) (int64, error) {
	oldList, _ := s.keyValues[key]
	updatedList, err := list.RightPushToOldList(values, oldList.data)
	if err != nil {
		return 0, err
	}

	s.keyValues[key] = NewEntry(updatedList)

	return int64(len(updatedList)), nil
}

func (s *InMemoryStore) ReadListRange(key string, fromIndex int, toIndex int) ([]string, error) {
	return list.ReadRangeFromStoreList(s.keyValues[key].data, fromIndex, toIndex)
}

func New() *InMemoryStore {
	return &InMemoryStore{
		keyValues: make(map[string]Entry),
	}
}
