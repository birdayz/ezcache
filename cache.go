package ezcache

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

type CacheConfig[K Key[K], V any] struct {
	capacity  int
	numShards int
	loader    LoaderFn[K, V]
}

func NewBuilder[K Key[K], V any]() *CacheConfig[K, V] {
	return &CacheConfig[K, V]{
		capacity:  1024,
		numShards: 1,
	}
}

func (cb *CacheConfig[K, V]) Capacity(capacity int) *CacheConfig[K, V] {
	cb.capacity = capacity
	return cb
}

func (cb *CacheConfig[K, V]) NumShards(numShards int) *CacheConfig[K, V] {
	cb.numShards = numShards
	return cb
}

func (cb *CacheConfig[K, V]) Loader(loader LoaderFn[K, V]) *CacheConfig[K, V] {
	cb.loader = loader
	return cb
}

func (cb *CacheConfig[K, V]) Build() *Cache[K, V] {
	return New[K, V](cb)
}

type HashCoder interface {
	HashCode() uint64
}

func New[K Key[K], V any](cfg *CacheConfig[K, V]) *Cache[K, V] {
	cache := Cache[K, V]{
		loaderFn:  cfg.loader,
		numShards: uint64(cfg.numShards),
		capacity:  cfg.capacity,
	}

	cache.shards = make([]*shard[K, V], 0, cache.numShards)
	for i := 0; i < int(cache.numShards); i++ {
		newShard := newShard[K, V]((cache.capacity / int(cache.numShards)) + 1)
		cache.shards = append(cache.shards, newShard)
	}

	return &cache
}

type LoaderFn[K Key[K], V any] func(key K) (value V, err error)

type Cache[K Key[K], V any] struct {
	loaderFn  LoaderFn[K, V]
	numShards uint64
	capacity  int

	shards []*shard[K, V]
}

func (c *Cache[K, V]) getShard(hash uint64) *shard[K, V] {
	return c.shards[hash%(c.numShards)]
}

func (c *Cache[K, V]) Set(key K, value V) {
	keyHash := key.HashCode()
	shard := c.getShard(keyHash)

	shard.set(key, keyHash, value)
}

func (c *Cache[K, V]) Get(key K) (V, error) {
	keyHash := key.HashCode()
	shard := c.getShard(keyHash)
	result, found := shard.get(key, keyHash)

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
		shard.set(key, keyHash, value)
		return value, nil
	}

	return result, nil
}

func (c *Cache[K, V]) Delete(key K) {
	keyHash := key.HashCode()
	shard := c.getShard(keyHash)

	shard.delete(key)
}
