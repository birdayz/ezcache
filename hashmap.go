package ezcache

import (
	"golang.org/x/exp/slices"
)

const loadFactor = 0.75

type Key[K any] interface {
	Equals(K) bool
	HashCoder
}

type bucket[K Key[K], V any] struct {
	slots []entry[K, V]
}

type entry[K Key[K], V any] struct {
	key   K
	value V
	hash  uint64
}

type HashMap[K Key[K], V any] struct {
	buckets []bucket[K, V]

	currentSize     int
	currentCapacity int
}

func NewHashMap[K Key[K], V any](initialCapacity int) *HashMap[K, V] {

	return &HashMap[K, V]{
		buckets:         make([]bucket[K, V], initialCapacity),
		currentSize:     0,
		currentCapacity: initialCapacity,
	}
}

func (h *HashMap[K, V]) maybeGrow() {
	if h.currentSize >= int(float64(h.currentCapacity)*loadFactor) {
		newBuckets := make([]bucket[K, V], h.currentCapacity*2)

		for i := range h.buckets {
			for d := range h.buckets[i].slots {
				hash := h.buckets[i].slots[d].hash
				idx := hash % uint64(len(newBuckets))

				newBuckets[idx].slots = append(newBuckets[idx].slots, h.buckets[i].slots[d])
			}
		}

		h.buckets = newBuckets

		h.currentCapacity = len(h.buckets)
	}

}

func (h *HashMap[K, V]) Set(key K, value V) bool {
	h.maybeGrow()

	hash := key.HashCode()
	bucket := &h.buckets[hash%uint64(len(h.buckets))]

	// Try to update
	for i := range bucket.slots {
		if bucket.slots[i].key.Equals(key) {
			bucket.slots[i].value = value
			return true
		}
	}

	// Insert new entry
	bucket.slots = append(bucket.slots, entry[K, V]{key, value, hash})
	h.currentSize++

	return false
}

func (h *HashMap[K, V]) Get(key K) (value V, found bool) {
	hash := key.HashCode()

	bucket := &h.buckets[hash%uint64(len(h.buckets))]

	for i := range bucket.slots {
		if bucket.slots[i].key.Equals(key) {
			return bucket.slots[i].value, true
		}
	}

	return *new(V), false
}

func (h *HashMap[K, V]) Delete(key K) (prev V, deleted bool) {
	hash := key.HashCode()

	bucket := &h.buckets[hash%uint64(len(h.buckets))]

	for i := range bucket.slots {
		if bucket.slots[i].key.Equals(key) {
			oldVal := bucket.slots[i].value
			bucket.slots = slices.Delete(bucket.slots, i, i+1)
			return oldVal, true
		}
	}
	return *new(V), false
}
