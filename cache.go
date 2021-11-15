package ezcache

import (
	"sync"

	"github.com/pkg/errors"
)

func New[K comparable, V comparable](loader LoaderFn[K, V], hasher HasherFn[K], numShards int) *Cache[K, V] {
	shards := make([]shard[K, V], 0, numShards)

	for i := 0; i < numShards; i++ {
		shards = append(shards, shard[K, V]{
			m:       sync.RWMutex{},
			buckets: make(map[uint64]bucket[K, V]),
			hasher:  hasher,
		})
	}

	return &Cache[K, V]{
		loaderFn: loader,
		hasherFn: hasher,
		shards:   shards,
	}
}

type LoaderFn[K comparable, V comparable] func(key K) (value V, err error)
type HasherFn[K comparable] func(key K) uint64

type Cache[K comparable, V comparable] struct {
	loaderFn LoaderFn[K, V]
	hasherFn HasherFn[K]
	shards   []shard[K, V]
}

func (c *Cache[K, V]) getShard(hash uint64) *shard[K, V] {
	return &c.shards[hash%(uint64(len(c.shards)))]
}

func (c *Cache[K, V]) Set(key K, value V) {
	keyHash := c.hasherFn(key)
	shard := c.getShard(keyHash)

	shard.set(key, value)
}

func (c *Cache[K, V]) Get(key K) (V, error) {
	keyHash := c.hasherFn(key)
	shard := c.getShard(keyHash)

	result, found := shard.get(key)

	if !found {
		value, err := c.loaderFn(key)
		if err != nil {
			return value, errors.Wrap(err, "failed to run loader")
		}

		// Since we don't hold the lock between get and set, it might be that we shadow other concurrent loads&writes.
		// It might be helpful to allow only one cache load per key concurrently, to avoid thundering herd etc.
		shard.set(key, value)
		return value, nil
	}

	return result, nil
}

func (c *Cache[K, V]) Delete(key K) {
	keyHash := c.hasherFn(key)
	shard := c.getShard(keyHash)

	shard.delete(key)
}
