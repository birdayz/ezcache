package ezcache

import (
	"sync"
	"time"
)

var (
	timeNow = time.Now
)

type shard[K interface {
	Equals(K) bool
	HashCoder
}, V comparable] struct {
	m sync.RWMutex

	buckets []bucketItem[K, V]

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
		buckets:    make([]bucketItem[K, V], 100000),
		linkedList: NewList[K](),
		capacity:   capacity,
		ttl:        time.Second * 2,
		ttls: NewHeap(func(t1, t2 *bucketItem[K, V]) int {
			if t1.expireAfter.After(t2.expireAfter) {
				return 1
			} else if t1.expireAfter.Before(t2.expireAfter) {
				return -1
			}
			return 0
		}),
	}
}

// set returns true if the value existed before
func (s *shard[K, V]) set(key K, keyHash uint64, value V) {
	s.m.Lock()
	defer s.m.Unlock()
	//s.clean()

	idx := int(keyHash) % len(s.buckets)
	if idx < 0 {
		idx = idx * -1
	}

	if !s.buckets[idx].filled {
		newElement := s.linkedList.PushFront(key)

		oldNext := s.buckets[idx].next

		s.buckets[idx] = bucketItem[K, V]{
			value:       value,
			key:         key,
			filled:      true,
			expireAfter: timeNow().Add(s.ttl),
			node:        newElement,
			heapElement: nil,
			next:        oldNext,
			previous:    nil,
		}

		newHeapItem := s.ttls.Push(&s.buckets[idx])
		s.buckets[idx].heapElement = newHeapItem

		return
	}

	// Try to find existing entry
	for pointer := &s.buckets[idx]; ; pointer = pointer.next {
		if pointer.key.Equals(key) {

			// Found, we can just replace
			pointer.value = value

			// Update expireAfter, "touch ttl/expiration"
			pointer.expireAfter = timeNow().Add(s.ttl)

			// The logic here is based on the ordering in the min-heap: we have
			// guaranteed ascending timestamp order. so we stop after seeing a
			// timestamp that's in the future. Without calling Fix, the now-increased
			// timestamp will stay at its location in the heap, which might be much
			// further in the front, and clean would stop when finding it.
			s.ttls.Fix(pointer.heapElement)

			// "Touch" it in LRU, set counts as "used"
			s.linkedList.MoveToFront(pointer.node)
			return
		}

		if pointer.next == nil {

			newElement := s.linkedList.PushFront(key)
			newItem := &bucketItem[K, V]{
				value:       value,
				key:         key,
				filled:      true,
				expireAfter: timeNow().Add(s.ttl),
				node:        newElement,
				heapElement: nil,
				next:        nil,
				previous:    nil,
			}
			newHeapItem := s.ttls.Push(newItem)
			newItem.heapElement = newHeapItem

			pointer.next = newItem
			return

		}
		// TODO, eviction not impl.

	}
}

func (s *shard[K, V]) clean() {
	for {
		if len(s.ttls.data) == 0 {
			return
		}

		item := s.ttls.Peek()
		if item.Item.expireAfter.Before(timeNow()) {
			// remove item
			res := s.delete(item.Item.key)
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

	//s.clean()

	idx := int(keyHash) % len(s.buckets)
	if idx < 0 {
		idx = idx * -1
	}

	for pointer := &s.buckets[idx]; pointer != nil; pointer = pointer.next {
		if pointer.key.Equals(key) && pointer.filled {
			s.linkedList.MoveToFront(pointer.node)
			return pointer.value, true
		}
	}

	return *new(V), false
}

func (s *shard[K, V]) Delete(key K) bool {
	panic("BIG DEL")
	s.m.Lock()
	defer s.m.Unlock()

	return s.delete(key)
}

func (s *shard[K, V]) delete(key K) bool {
	keyHash := key.HashCode()
	idx := int(keyHash) % len(s.buckets)
	if idx < 0 {
		idx = idx * -1
	}

	if !s.buckets[idx].filled && s.buckets[idx].next == nil {
		return false
	}

	for pointer := &s.buckets[idx]; pointer != nil; pointer = pointer.next {
		if pointer.key.Equals(key) {
			pointer.filled = false
			if pointer.previous != nil {
				pointer.previous.next = nil
			}

			if pointer.next != nil {
				pointer.next.previous = pointer.previous
			}

			s.linkedList.Remove(pointer.node)
			s.ttls.Remove(pointer.heapElement)

			return true

		}
	}

	// if s.buckets[idx].filled {
	// 	//for i, bi := range bucket.items {
	// 	if s.buckets[idx].key.Equals(key) {
	// 		// We actually have the key
	//
	// 		s.linkedList.Remove(s.buckets[idx].node)
	// 		// Can probably be optimized
	// 		//bucket.items = slices.Delete(bucket.items, i, i+1)
	// 		s.buckets[idx].filled = false
	// 		s.ttls.Remove(s.buckets[idx].heapElement)
	//
	// 		return true
	// 	}
	// 	//}
	// }
	return false
}

// bucket

type bucketItem[K interface {
	HashCoder
	Equals(K) bool
}, V comparable] struct {
	value V
	key   K

	filled bool

	expireAfter time.Time

	// LinkedList node pointer, used for LRU eviction
	node *Element[K]

	// Pointer to heap item, used for TTL
	heapElement *HeapElement[*bucketItem[K, V]]

	// Pointers to next/prev items in bucket
	next     *bucketItem[K, V]
	previous *bucketItem[K, V]
}
