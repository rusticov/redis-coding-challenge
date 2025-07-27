package store

import (
	"redis-challenge/internal/store/list"
	"strconv"
)

const maximumTimeInFuture = int64(9223372036854775807)

type entry struct {
	data                     any
	expiryTimeInMilliseconds int64
}

type InMemoryStore struct {
	keyEntries map[string]entry
	clock      Clock
}

func (s *InMemoryStore) Exists(key string) bool {
	_, ok := s.readEntry(key)
	return ok
}

func (s *InMemoryStore) ReadString(key string) (string, error) {
	if e, ok := s.readEntry(key); ok {
		if text, ok := e.data.(string); ok {
			return text, nil
		}
		return "", ErrorWrongOperationType
	}
	return "", ErrorKeyNotFound
}

func (s *InMemoryStore) readEntry(key string) (entry, bool) {
	if keyEntry, ok := s.keyEntries[key]; ok {
		expirationTime := keyEntry.expiryTimeInMilliseconds

		if expirationTime > s.clock() {
			return keyEntry, true
		} else {
			delete(s.keyEntries, key)
		}
	}
	return entry{}, false
}

func (s *InMemoryStore) Delete(key string) bool {
	existed := s.Exists(key)
	delete(s.keyEntries, key)
	return existed
}

func (s *InMemoryStore) Increment(key string, incrementBy int64) (int64, error) {
	value, err := s.readInteger(key)
	if err != nil {
		return 0, err
	}
	value += incrementBy

	stringValue := strconv.FormatInt(value, 10)

	s.Write(key, stringValue, ExpiryOptionNone, 0)

	return value, nil
}

func (s *InMemoryStore) readInteger(key string) (int64, error) {
	if e, hasEntry := s.readEntry(key); hasEntry {
		if text, ok := e.data.(string); ok {
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
	oldList, _ := s.readEntry(key)
	updatedList, ok := list.LeftPushToOldList(values, oldList.data)
	if !ok {
		return 0, ErrorWrongOperationType
	}

	s.keyEntries[key] = entry{
		data:                     updatedList,
		expiryTimeInMilliseconds: maximumTimeInFuture,
	}

	return int64(len(updatedList)), nil
}

func (s *InMemoryStore) RightPush(key string, values []string) (int64, error) {
	oldList, _ := s.readEntry(key)
	updatedList, err := list.RightPushToOldList(values, oldList.data)
	if err != nil {
		return 0, ErrorWrongOperationType
	}

	s.keyEntries[key] = entry{
		data:                     updatedList,
		expiryTimeInMilliseconds: maximumTimeInFuture,
	}

	return int64(len(updatedList)), nil
}

func (s *InMemoryStore) ReadListRange(key string, fromIndex int, toIndex int) ([]string, error) {
	if values, ok := list.ReadRangeFromStoreList(s.keyEntries[key].data, fromIndex, toIndex); ok {
		return values, nil
	}
	return nil, ErrorWrongOperationType
}

func (s *InMemoryStore) Write(key string, value string, expiryOption ExpiryOption, expiry int64) {
	s.keyEntries[key] = entry{
		data:                     value,
		expiryTimeInMilliseconds: s.expiryTimeInMilliseconds(key, expiryOption, expiry),
	}
}

func (s *InMemoryStore) expiryTimeInMilliseconds(key string, expiryOption ExpiryOption, expiry int64) int64 {
	switch expiryOption {
	case ExpiryOptionNone:
		return maximumTimeInFuture

	case ExpiryOptionExpiryMilliseconds:
		return s.clock() + expiry
	case ExpiryOptionExpirySeconds:
		return s.clock() + expiry*1000

	case ExpiryOptionExpiryUnixTimeInMilliseconds:
		return expiry
	case ExpiryOptionExpiryUnixTimeInSeconds:
		return expiry * 1000

	case ExpiryOptionExpiryKeepTTL:
		if keyEntry, ok := s.keyEntries[key]; ok {
			return keyEntry.expiryTimeInMilliseconds
		}
	}

	return maximumTimeInFuture
}

func (s *InMemoryStore) Size() int {
	return len(s.keyEntries)
}

func New() *InMemoryStore {
	return NewWithClock(CurrentSystemTime)
}

func NewWithClock(clock Clock) *InMemoryStore {
	return &InMemoryStore{
		keyEntries: make(map[string]entry),
		clock:      clock,
	}
}
