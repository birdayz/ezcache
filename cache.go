package ezcache

import (
	"errors"
	"fmt"
	"sync"
)

var ErrNotFound = errors.New("not found")

type HashCoder interface {
	HashCode() uint64
}

func New[K interface {
	HashCoder
	Equals(K) bool
}, V comparable](loader LoaderFn[K, V], numShards int, capacity int) *Cache[K, V] {
	shards := make([]shard[K, V], 0, numShards)

	for i := 0; i < numShards; i++ {
		shards = append(shards, shard[K, V]{
			m:          sync.RWMutex{},
			buckets:    make(map[uint64]*bucket[K, V]),
			linkedList: NewList[K](),
			capacity:   (capacity / numShards) + 1,
		})
	}

	return &Cache[K, V]{
		loaderFn: loader,
		shards:   shards,
	}
}

type LoaderFn[K any, V any] func(key K) (value V, err error)

type Cache[K interface {
	HashCoder
	Equals(K) bool
}, V comparable] struct {
	loaderFn LoaderFn[K, V]
	shards   []shard[K, V]
}

func (c *Cache[K, V]) getShard(hash uint64) *shard[K, V] {
	return &c.shards[hash%(uint64(len(c.shards)))]
}

func (c *Cache[K, V]) Set(key K, value V) {
	keyHash := key.HashCode()
	shard := c.getShard(keyHash)

	shard.set(key, value)
}

func (c *Cache[K, V]) Get(key K) (V, error) {
	keyHash := key.HashCode()
	shard := c.getShard(keyHash)

	result, found := shard.get(key)

	if !found {
		if c.loaderFn == nil {
			return *new(V), ErrNotFound
		}
		value, err := c.loaderFn(key)
		if err != nil {
			return *new(V), fmt.Errorf("failed to run loader: %w", err)
		}

		// Since we don't hold the lock between get and set, it might be that we shadow other concurrent loads&writes.
		// It might be helpful to allow only one cache load per key concurrently, to avoid thundering herd etc.
		shard.set(key, value)
		return value, nil
	}

	return result, nil
}

func (c *Cache[K, V]) Delete(key K) {
	keyHash := key.HashCode()
	shard := c.getShard(keyHash)

	shard.delete(key)
}
