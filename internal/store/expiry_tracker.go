package store

import "math/rand"

type ExpiryTracker struct {
	keys     []string
	keyIsSet map[string]struct{}
}

func (t *ExpiryTracker) SelectKeys(count int) []string {
	totalKeyCount := len(t.keys)

	if count >= totalKeyCount {
		return t.keys
	}

	indexes := make(map[int]struct{}, count)

	var keysToReturn = make([]string, count)
	for i := range count {
		index := rand.Intn(totalKeyCount)
		for {
			if _, ok := indexes[index]; !ok {
				indexes[index] = struct{}{}
				break
			}

			index++
			if index == totalKeyCount {
				index = 0
			}
		}

		keysToReturn[i] = t.keys[index]
	}

	return keysToReturn
}

func (t *ExpiryTracker) AddKey(key string) {
	if t != nil {
		if _, ok := t.keyIsSet[key]; !ok {
			t.keyIsSet[key] = struct{}{}
			t.keys = append(t.keys, key)
		}
	}
}

func (t *ExpiryTracker) RemoveKey(key string) {
	if t != nil {
		delete(t.keyIsSet, key)
		for i, k := range t.keys {
			if k == key {
				t.keys = append(t.keys[:i], t.keys[i+1:]...)
				break
			}
		}
	}
}

func NewExpiryTracker() *ExpiryTracker {
	return &ExpiryTracker{
		keyIsSet: make(map[string]struct{}),
	}
}
