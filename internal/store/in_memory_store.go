package store

import (
	"redis-challenge/internal/store/list"
	"strconv"
)

type entry struct {
	data                     any
	expiryTimeInMilliseconds int64
}

func newEntry(data any) entry {
	return entry{data: data}
}

func (e entry) Data() any {
	return e.data
}

type InMemoryStore struct {
	keyEntries map[string]entry
	clock      Clock
}

func (s *InMemoryStore) Exists(key string) bool {
	_, ok := s.keyEntries[key]
	return ok
}

func (s *InMemoryStore) Delete(key string) bool {
	existed := s.Exists(key)
	delete(s.keyEntries, key)
	return existed
}

func (s *InMemoryStore) Write(key string, value string) {
	s.keyEntries[key] = newEntry(value)
}

func (s *InMemoryStore) ReadString(key string) (string, error) {
	if keyEntry, ok := s.keyEntries[key]; ok {
		expirationTime := keyEntry.expiryTimeInMilliseconds

		if expirationTime == 0 || expirationTime > s.clock() {
			return keyEntry.data.(string), nil
		}
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
	s.keyEntries[key] = newEntry(stringValue)

	return value, nil
}

func (s *InMemoryStore) readInteger(key string) (int64, error) {
	if entry, ok := s.keyEntries[key]; ok {
		if text, ok := entry.data.(string); ok {
			value, err := strconv.ParseInt(text, 10, 64)
			if err != nil {
				return 0, ErrorNotAnInteger
			}
			return value, nil
		}
		return 0, ErrorWrongOperationType
	}
	return 0, nil
}

func (s *InMemoryStore) LeftPush(key string, values []string) (int64, error) {
	oldList, _ := s.keyEntries[key]
	updatedList, err := list.LeftPushToOldList(values, oldList.data)
	if err != nil {
		return 0, err
	}

	s.keyEntries[key] = newEntry(updatedList)

	return int64(len(updatedList)), nil
}

func (s *InMemoryStore) RightPush(key string, values []string) (int64, error) {
	oldList, _ := s.keyEntries[key]
	updatedList, err := list.RightPushToOldList(values, oldList.data)
	if err != nil {
		return 0, err
	}

	s.keyEntries[key] = newEntry(updatedList)

	return int64(len(updatedList)), nil
}

func (s *InMemoryStore) ReadListRange(key string, fromIndex int, toIndex int) ([]string, error) {
	return list.ReadRangeFromStoreList(s.keyEntries[key].data, fromIndex, toIndex)
}

func (s *InMemoryStore) WriteWithExpiry(key string, value string, expiryOption ExpiryOption, expiry int64) {
	s.keyEntries[key] = entry{
		data:                     value,
		expiryTimeInMilliseconds: expiry,
	}
}

func New() *InMemoryStore {
	return NewWithCLock(CurrentSystemTime)
}

func NewWithCLock(clock Clock) *InMemoryStore {
	return &InMemoryStore{
		keyEntries: make(map[string]entry),
		clock:      clock,
	}
}
