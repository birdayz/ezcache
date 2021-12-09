package ezcache

import (
	"fmt"

	"golang.org/x/exp/slices"
)

type bucket[K interface {
	HashCoder
	Equals(K) bool
}, V any] struct {
	slots []slot[K, V]
}

type slot[K interface {
	HashCoder
	Equals(K) bool
}, V any] struct {
	key   K
	value V
}

type HashMap[K interface {
	HashCoder
	Equals(K) bool
}, V any] struct {
	buckets []bucket[K, V]

	currentSize int

	capacity int
}

func NewHashMap[K interface {
	HashCoder
	Equals(K) bool
}, V any]() *HashMap[K, V] {

	return &HashMap[K, V]{
		buckets:     make([]bucket[K, V], 16),
		currentSize: 0,
		capacity:    16,
	}
}

// insert inserts a new entry into the bucket
func (h *HashMap[K, V]) insert(bucket *bucket[K, V], key K, value V) {

	// Check if we need to grow
	loadFactor := 0.75

	if h.currentSize >= int(float64(h.capacity)*loadFactor) {
		fmt.Println("Increase", loadFactor, h.currentSize, h.capacity)

		// need to grow
		slices.Grow(h.buckets, h.capacity)
		h.capacity = h.capacity * 2
	}

	bucket.slots = append(bucket.slots, slot[K, V]{key, value})
	h.currentSize++
}

func (h *HashMap[K, V]) Set(key K, value V) bool {
	hash := key.HashCode()

	bucket := &h.buckets[hash%uint64(len(h.buckets))]

	for i := range bucket.slots {
		if bucket.slots[i].key.Equals(key) {
			bucket.slots[i].value = value
			return true
		}
	}

	h.insert(bucket, key, value)

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
