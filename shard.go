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
}, V any] struct {
	m sync.RWMutex

	dataMap *HashMap[K, *bucketItem[K, V]]

	linkedList *List[K]
	capacity   int
	ttl        time.Duration

	ttls *Heap[*bucketItem[K, V]]
}

func newShard[K interface {
	Equals(K) bool
	HashCoder
}, V any](capacity int) *shard[K, V] {

	var initialMapCapacity = capacity
	if initialMapCapacity < 16 {
		initialMapCapacity = 16
	}

	return &shard[K, V]{
		m:          sync.RWMutex{},
		dataMap:    NewHashMap[K, *bucketItem[K, V]](initialMapCapacity),
		linkedList: NewList[K](),
		capacity:   capacity,
		ttl:        time.Second * 50,
		ttls: NewHeap(func(t1, t2 *bucketItem[K, V]) int {
			if t1.expireAt > t2.expireAt {
				return 1
			} else if t1.expireAt < t2.expireAt {
				return -1
			}

			return 0
		}, capacity),
	}
}

// set returns true if the value existed before
func (s *shard[K, V]) set(key K, keyHash uint64, value V) {
	s.m.Lock()
	defer s.m.Unlock()

	s.clean()

	// This could be optimized with a very specific call that does the get and
	// update at once
	entry, ok := s.dataMap.Get(key)
	if !ok {

		if s.linkedList.Len() >= s.capacity {
			keyToRemove := s.linkedList.Back()
			s.delete(keyToRemove.Value)
		}

		// Not found
		newElement := s.linkedList.PushFront(key)

		newItem := bucketItem[K, V]{
			value:       value,
			expireAt:    timeNow().Add(s.ttl).UnixMilli(),
			node:        newElement,
			heapElement: nil,
		}

		newHeapItem := s.ttls.Push(&newItem)
		newItem.heapElement = newHeapItem
		s.dataMap.Set(key, &newItem)

		return

	} else {
		entry.expireAt = timeNow().Add(s.ttl).UnixMilli() // TODO: store ttls somewhere else, not in the map entry
		entry.value = value
		s.ttls.Fix(entry.heapElement)
		s.linkedList.MoveToFront(entry.node)
		s.dataMap.Set(key, entry)
	}
}

func (s *shard[K, V]) clean() {
	for {
		if len(s.ttls.data) == 0 {
			return
		}

		item := s.ttls.Peek()
		if item.Item.expireAt <= timeNow().UnixMilli() {
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

func (s *shard[K, V]) get(key K, keyHash uint64) (V, bool) {
	s.m.Lock()
	defer s.m.Unlock()

	s.clean()

	data, ok := s.dataMap.Get(key)
	if !ok {
		return *new(V), false
	}

	s.linkedList.MoveToFront(data.node)
	return data.value, true

}

func (s *shard[K, V]) Delete(key K) bool {
	s.m.Lock()
	defer s.m.Unlock()

	return s.delete(key)
}

func (s *shard[K, V]) delete(key K) bool {
	oldVal, deleted := s.dataMap.Delete(key)

	if deleted {
		s.linkedList.Remove(oldVal.node)
		s.ttls.Remove(oldVal.heapElement)
		return true
	}

	return false
}

// bucket

type bucketItem[K interface {
	HashCoder
	Equals(K) bool
}, V any] struct {
	value V

	expireAt int64 // exact timestamp, at which the entry is considered expired

	// LinkedList node pointer, used for LRU eviction
	node *Element[K]

	// Pointer to heap item, used for TTL
	heapElement *HeapElement[*bucketItem[K, V]]
}
