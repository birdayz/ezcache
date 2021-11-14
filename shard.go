package ezcache

import "sync"

type shard[K comparable, V comparable] struct {
	m       sync.RWMutex
	buckets map[uint64]bucket[K, V]

	hasher HasherFn[K]
}

func (s *shard[K, V]) set(key K, value V) {
	s.m.Lock()
	defer s.m.Unlock()

	keyHash := s.hasher(key)

	b, found := s.buckets[keyHash]
	if !found {
		newBucket := bucket[K, V]{
			items: make(map[K]V),
		}

		s.buckets[keyHash] = newBucket
		b = newBucket

	}

	b.items[key] = value
}

func (s *shard[K, V]) get(key K) V {
	s.m.RLock()
	defer s.m.RUnlock()

	keyHash := s.hasher(key)

	if bucket, found := s.buckets[keyHash]; found {
		if item, found := bucket.items[key]; found {
			return item
		}
	}

	return *new(V)
}

func (s *shard[K, V]) delete(key K) bool {
	s.m.Lock()
	defer s.m.Unlock()

	keyHash := s.hasher(key)

	b, found := s.buckets[keyHash]
	if found {
		if _, found := b.items[key]; found {
			delete(b.items, key)
			return true
		}

	}

	return false
}

// bucket

type bucket[K comparable, V comparable] struct {
	items map[K]V
}
