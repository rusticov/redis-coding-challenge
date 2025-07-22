package command_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"redis-challenge/internal/command"
	"redis-challenge/internal/protocol"
	"redis-challenge/internal/store"
	"sync"
	"testing"
)

func TestConcurrentAdding(t *testing.T) {

	t.Run("first write wins when there are competing SET commands with NX", func(t *testing.T) {
		dataStore := store.New()

		cmd1, errorData := command.Validate(protocol.Array{Data: []protocol.Data{
			protocol.NewBulkString("SET"),
			protocol.NewBulkString("key"),
			protocol.NewBulkString("value 1"),
			protocol.NewBulkString("NX"),
		}})
		require.Nil(t, errorData)

		cmd2, errorData := command.Validate(protocol.Array{Data: []protocol.Data{
			protocol.NewBulkString("SET"),
			protocol.NewBulkString("key"),
			protocol.NewBulkString("value 2"),
			protocol.NewBulkString("NX"),
		}})
		require.Nil(t, errorData)

		firstCommandWait := sync.WaitGroup{}
		firstCommandWait.Add(1)

		secondCommandWait := sync.WaitGroup{}
		secondCommandWait.Add(1)

		var wg sync.WaitGroup
		wg.Add(2)

		var value1 protocol.Data
		var value2 protocol.Data

		go func() {
			defer wg.Done()

			var err error
			value1, err = cmd1.Execute(&waitingStore{
				store:       dataStore,
				beforeWrite: &secondCommandWait,
				waitToWrite: &firstCommandWait,
			})
			require.NoError(t, err)
		}()

		go func() {
			defer wg.Done()

			var err error
			value2, err = cmd2.Execute(&waitingStore{
				store:          dataStore,
				waitBeforeRead: &secondCommandWait,
				afterWriting:   &firstCommandWait,
			})
			require.NoError(t, err)
		}()

		wg.Wait()

		assert.Nil(t, value1, "first command does not set value because second command writes first")
		assert.Equal(t, protocol.NewSimpleString("OK"), value2,
			"second command should write first")
	})
}

type waitingStore struct {
	store          store.Store
	waitBeforeRead *sync.WaitGroup
	beforeWrite    *sync.WaitGroup
	waitToWrite    *sync.WaitGroup
	afterWriting   *sync.WaitGroup
}

func (s *waitingStore) CompareAndSwap(key string, oldValue, newValue string) (swapped bool) {
	if s.beforeWrite != nil {
		s.beforeWrite.Done()
	}
	if s.waitToWrite != nil {
		s.waitToWrite.Wait()
	}

	swapped = s.store.CompareAndSwap(key, oldValue, newValue)

	if s.afterWriting != nil {
		s.afterWriting.Done()
	}

	return
}

func (s *waitingStore) LoadOrStore(key string, defaultValue string) (string, bool) {
	if s.beforeWrite != nil {
		s.beforeWrite.Done()
	}
	if s.waitToWrite != nil {
		s.waitToWrite.Wait()
	}

	value, loaded := s.store.LoadOrStore(key, defaultValue)

	if s.afterWriting != nil {
		s.afterWriting.Done()
	}

	return value, loaded
}

func (s *waitingStore) Get(key string) (string, bool) {
	if s.waitBeforeRead != nil {
		s.waitBeforeRead.Wait()
	}
	return s.store.Get(key)
}
