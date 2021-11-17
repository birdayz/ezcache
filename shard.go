package ezcache

import "sync"

type shard[K interface {
	Equals(K) bool
	HashCode() int
}, V comparable] struct {
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
			items: make([]bucketItem[K, V], 0, 0),
		}

		s.buckets[keyHash] = newBucket
		b = newBucket

	}

	// Try to find entry for key
	for _, buckItem := range b.items {
		if buckItem.key.Equals(key) {
			buckItem.value = value
			return
		}
	}

	b.items = append(b.items, bucketItem[K, V]{
		key:   key,
		value: value,
	})
}

func (s *shard[K, V]) get(key K) (V, bool) {
	s.m.RLock()
	defer s.m.RUnlock()

	keyHash := s.hasher(key)

	if bucket, found := s.buckets[keyHash]; found {
		for _, bucketItem := range bucket.items {
			if bucketItem.key.Equals(key) {
				return bucketItem.value, true
			}
		}
	}

	return *new(V), false
}

func (s *shard[K, V]) delete(key K) bool {
	s.m.Lock()
	defer s.m.Unlock()

	keyHash := s.hasher(key)

	if bucket, found := s.buckets[keyHash]; found {
		for i, bucketItem := range bucket.items {
			if bucketItem.key.Equals(key) {
				_ = i
				//bucket.items[len(bucket.items)-1], bucket.items[i] = new(bucketItem[K, V]), bucket.items[len(bucket.items)-1]
				return true
			}
		}
	}

	return false
}

// bucket

type bucketItem[K interface {
	Equals(K) bool
	HashCode() int
}, V comparable] struct {
	key   K
	value V
}

type bucket[K interface {
	Equals(K) bool
	HashCode() int
}, V comparable] struct {
	items []bucketItem[K, V]
}
