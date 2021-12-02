package ezcache

import (
	"sync"

	"golang.org/x/exp/slices"
)

type shard[K interface {
	Equals(K) bool
	HashCoder
}, V comparable] struct {
	m       sync.RWMutex
	buckets map[uint64]*bucket[K, V]

	linkedList *List[K]
	capacity   int
}

func newShard[K interface {
	Equals(K) bool
	HashCoder
}, V comparable](capacity int) *shard[K, V] {

	return &shard[K, V]{
		m:          sync.RWMutex{},
		buckets:    map[uint64]*bucket[K, V]{},
		linkedList: NewList[K](),
		capacity:   capacity,
	}
}

// set returns true if the value existed before
func (s *shard[K, V]) set(key K, keyHash uint64, value V) {
	s.m.Lock()
	defer s.m.Unlock()

	b, found := s.buckets[keyHash]
	if !found {
		newBucket := bucket[K, V]{
			items: make([]*bucketItem[K, V], 0, 4),
		}

		s.buckets[keyHash] = &newBucket
		b = &newBucket

	}

	// Try to find entry for key
	for _, buckItem := range b.items {
		if buckItem.node.Value.Equals(key) {

			// Found, we can just replace
			buckItem.value = value

			// "Touch" it in LRU, set counts as "used"
			s.linkedList.MoveToFront(buckItem.node)
			return
		}
	}

	// Key is not in here yet, we have to add it

	// Check if we need to evict
	if s.linkedList.Len() >= s.capacity {

		keyToRemove := s.linkedList.Back()

		// Remove item. delete also removed the LRU entry, no need to specifically
		// do it here
		s.delete(keyToRemove.Value)
	}

	newElement := s.linkedList.PushFront(key)

	b.items = append(b.items, &bucketItem[K, V]{
		value: value,
		node:  newElement,
	})
}

func (s *shard[K, V]) get(key K, keyHash uint64) (V, bool) {
	s.m.Lock()
	defer s.m.Unlock()

	// TODO, try to use RLock - it's not simple, because this func actually does
	// modify : the linkedList.

	if bucket, found := s.buckets[keyHash]; found {
		for _, bucketItem := range bucket.items {
			if bucketItem.node.Value.Equals(key) {
				s.linkedList.MoveToFront(bucketItem.node)
				return bucketItem.value, true
			}
		}
	}

	return *new(V), false
}

func (s *shard[K, V]) Delete(key K) bool {
	s.m.Lock()
	defer s.m.Unlock()

	return s.delete(key)
}

func (s *shard[K, V]) delete(key K) bool {
	keyHash := key.HashCode()

	if bucket, found := s.buckets[keyHash]; found {
		for i, bi := range bucket.items {
			if bi.node.Value.Equals(key) {
				// We actually have the key

				s.linkedList.Remove(bi.node)
				// Can probably be optimized
				bucket.items = slices.Delete(bucket.items, i, i+1)
				return true
			}
		}
	}

	return false
}

// bucket

type bucketItem[K interface {
	HashCoder
	Equals(K) bool
}, V comparable] struct {
	value V
	node  *Element[K]
}

type bucket[K interface {
	HashCoder
	Equals(K) bool
}, V comparable] struct {
	items []*bucketItem[K, V]
}
