package ezcache

import (
	"sync"
)

func New[K comparable, V comparable](loader LoaderFn[K, V], hasher HasherFn[K]) *Cache[K, V] {
	return &Cache[K, V]{
		m: sync.Mutex{},

		tree: make(map[K]V),

		loaderFn: loader,
		hasherFn: hasher,
	}
}

type LoaderFn[K comparable, V comparable] func(key K) (value V, err error)

type HasherFn[K comparable] func(key K) uint32

type Cache[K comparable, V comparable] struct {
	m    sync.Mutex
	tree map[K]V

	loaderFn LoaderFn[K, V]
	hasherFn HasherFn[K]
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.m.Lock()
	defer c.m.Unlock()

	c.tree[key] = value
}

func (c *Cache[K, V]) Get(key K) (V, error) {
	c.m.Lock()
	defer c.m.Unlock()

	result, ok := c.tree[key]
	if ok {
		return result, nil
	}

	val, err := c.loaderFn(key)
	return val, err

}
