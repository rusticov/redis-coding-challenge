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
	keyEntries    map[string]entry
	clock         Clock
	expiryTracker *ExpiryTracker
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

		if expirationTime > s.clock.Now() {
			return keyEntry, true
		} else {
			s.expiryTracker.RemoveKey(key)
			delete(s.keyEntries, key)
		}
	}
	return entry{}, false
}

func (s *InMemoryStore) Delete(key string) bool {
	existed := s.Exists(key)

	delete(s.keyEntries, key)
	s.expiryTracker.RemoveKey(key)

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
	updatedList, ok := list.RightPushToOldList(values, oldList.data)
	if !ok {
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
	if expiryOption == ExpiryOptionNone {
		s.expiryTracker.RemoveKey(key)
	} else {
		s.expiryTracker.AddKey(key)
	}

	expiryTimestamp, ok := s.expiryTimeInMilliseconds(key, expiryOption, expiry)

	if ok {
		s.keyEntries[key] = entry{
			data:                     value,
			expiryTimeInMilliseconds: expiryTimestamp,
		}
	}
}

func (s *InMemoryStore) expiryTimeInMilliseconds(key string, expiryOption ExpiryOption, expiry int64) (int64, bool) {
	now := s.clock.Now()

	var expiryTimestamp int64
	switch expiryOption {
	case ExpiryOptionNone:
		expiryTimestamp = maximumTimeInFuture

	case ExpiryOptionExpiryMilliseconds:
		expiryTimestamp = now + expiry
	case ExpiryOptionExpirySeconds:
		expiryTimestamp = now + expiry*1000

	case ExpiryOptionExpiryUnixTimeInMilliseconds:
		expiryTimestamp = expiry
	case ExpiryOptionExpiryUnixTimeInSeconds:
		expiryTimestamp = expiry * 1000

	case ExpiryOptionExpiryKeepTTL:
		if keyEntry, ok := s.keyEntries[key]; ok {
			expiryTimestamp = keyEntry.expiryTimeInMilliseconds
		} else {
			expiryTimestamp = maximumTimeInFuture
		}
	}

	return expiryTimestamp, expiryTimestamp > now
}

func (s *InMemoryStore) Size() int {
	return len(s.keyEntries)
}

func (s *InMemoryStore) WithExpiryTracker(tracker *ExpiryTracker) *InMemoryStore {
	s.expiryTracker = tracker
	return s
}

func New() *InMemoryStore {
	return NewWithClock(SystemClock{})
}

func NewWithClock(clock Clock) *InMemoryStore {
	return &InMemoryStore{
		keyEntries: make(map[string]entry),
		clock:      clock,
	}
}
