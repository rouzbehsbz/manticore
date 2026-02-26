package util

import "sync"

type SyncMap[K comparable, V any] struct {
	unsafeMap map[K]V
	mu        sync.RWMutex
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		unsafeMap: make(map[K]V),
		mu:        sync.RWMutex{},
	}
}

func (s *SyncMap[K, V]) Insert(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.unsafeMap[key] = value
}

func (s *SyncMap[K, V]) Get(key K) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.unsafeMap[key]
	return value, ok
}

func (s *SyncMap[K, V]) Delete(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.unsafeMap, key)
}
