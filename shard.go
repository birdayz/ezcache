package ezcache

import (
	"sync"
	"time"

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
	ttl        time.Duration

	ttls *Heap[*bucketItem[K, V]]
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
		ttls: NewHeap(func(t1, t2 *bucketItem[K, V]) int {
			if t1.expireAfter.After(t2.expireAfter) {
				return 1
			} else if t1.expireAfter.Before(t2.expireAfter) {
				return -1
			}
			return 0
		}),
		ttl: time.Second * 2, //TODO
	}
}

// set returns true if the value existed before
func (s *shard[K, V]) set(key K, keyHash uint64, value V) {
	s.m.Lock()
	defer s.m.Unlock()
	s.clean()

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

			// Update expireAfter, "touch ttl/expiration"
			buckItem.expireAfter = time.Now().Add(s.ttl)

			// The logic here is based on the ordering in the min-heap: we have
			// guaranteed ascending timestamp order. so we stop after seeing a
			// timestamp that's in the future. Without calling Fix, the now-increased
			// timestamp will stay at its location in the heap, which might be much
			// further in the front, and clean would stop when finding it.
			s.ttls.Fix(buckItem.heapElement)

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

	bi := &bucketItem[K, V]{
		value:       value,
		node:        newElement,
		expireAfter: time.Now().Add(s.ttl),
	}

	newHeapItem := s.ttls.Push(bi)
	bi.heapElement = newHeapItem
	b.items = append(b.items, bi)
}

func (s *shard[K, V]) clean() {
	for {
		if len(s.ttls.data) == 0 {
			return
		}

		item := s.ttls.Peek()
		if item.Item.expireAfter.Before(time.Now()) {
			// remove item
			res := s.delete(item.Item.node.Value)
			if !res {
				panic("bug - delete was unsuccessful. This means that fundamental invariants are broken, and the cache's internal state is most likely not consistent anymore")
			}
		} else {
			return
		}
	}
}

// Cache ttl: does a get extend retention?
func (s *shard[K, V]) get(key K, keyHash uint64) (V, bool) {
	s.m.Lock()
	defer s.m.Unlock()

	s.clean()

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
				s.ttls.Remove(bi.heapElement)

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
	value       V
	node        *Element[K]
	expireAfter time.Time

	// Pointer to heap item
	heapElement *HeapElement[*bucketItem[K, V]]
}

type bucket[K interface {
	HashCoder
	Equals(K) bool
}, V comparable] struct {
	items []*bucketItem[K, V]
}
