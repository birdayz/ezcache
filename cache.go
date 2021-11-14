package ezcache

import (
	"sync"
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
		m:        sync.Mutex{},
		loaderFn: loader,
		hasherFn: hasher,
		shards:   shards,
	}
}

type LoaderFn[K comparable, V comparable] func(key K) (value V, err error)

type HasherFn[K comparable] func(key K) uint64

type Cache[K comparable, V comparable] struct {
	m sync.Mutex

	loaderFn LoaderFn[K, V]
	hasherFn HasherFn[K]

	shards []shard[K, V]
}

func (c *Cache[K, V]) getShard(hash uint64) *shard[K, V] {
	return &c.shards[hash%(uint64(len(c.shards)))]
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.m.Lock()
	defer c.m.Unlock()

	keyHash := c.hasherFn(key)
	shard := c.getShard(keyHash)

	shard.set(key, value)
}

func (c *Cache[K, V]) Get(key K) (V, error) {
	c.m.Lock()
	defer c.m.Unlock()

	keyHash := c.hasherFn(key)
	shard := c.getShard(keyHash)

	return shard.get(key), nil
}
