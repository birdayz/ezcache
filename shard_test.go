package ezcache

import (
	"sync"
	"testing"

	"gotest.tools/v3/assert"
)

func TestGet(t *testing.T) {
	shard := &shard[StringKey, string]{
		m:          sync.RWMutex{},
		buckets:    map[uint64]*bucket[StringKey, string]{},
		linkedList: NewList[StringKey](),
		capacity:   10,
	}

	shard.set("abc", "def")
	res, ok := shard.get("abc")

	assert.Equal(t, ok, true)
	assert.Equal(t, res, "def")
}

func TestSetGetEvict(t *testing.T) {
	shard := &shard[StringKey, string]{
		m:          sync.RWMutex{},
		buckets:    map[uint64]*bucket[StringKey, string]{},
		linkedList: NewList[StringKey](),
		capacity:   2,
	}

	shard.set("first", "def")
	res, ok := shard.get("first")
	assert.Equal(t, ok, true)
	assert.Equal(t, res, "def")

	shard.set("second", "second_val")
	res, ok = shard.get("second")
	assert.Equal(t, ok, true)
	assert.Equal(t, res, "second_val")

	// Trigger eviction
	shard.set("third", "third_val")
	res, ok = shard.get("third")
	assert.Equal(t, ok, true)
	assert.Equal(t, res, "third_val")

	res, ok = shard.get("first")
	assert.Equal(t, ok, false)
}

func TestSetGetEvictOrder(t *testing.T) {
	shard := &shard[StringKey, string]{
		m:          sync.RWMutex{},
		buckets:    map[uint64]*bucket[StringKey, string]{},
		linkedList: NewList[StringKey](),
		capacity:   2,
	}

	shard.set("first", "def")
	res, ok := shard.get("first")
	assert.Equal(t, ok, true)
	assert.Equal(t, res, "def")

	shard.set("second", "second_val")
	res, ok = shard.get("second")
	assert.Equal(t, ok, true)
	assert.Equal(t, res, "second_val")

	// Issue get, so first is recently used; shall not be evicted
	shard.get("first")

	// Trigger eviction
	shard.set("third", "third_val")
	res, ok = shard.get("third")
	assert.Equal(t, ok, true)
	assert.Equal(t, res, "third_val")

	// First has to be preserved, because we touched it
	res, ok = shard.get("first")
	assert.Equal(t, ok, true)

	res, ok = shard.get("second")
	assert.Equal(t, ok, false)
}

func TestGetDoesNotExist(t *testing.T) {
	shard := &shard[StringKey, string]{
		m:       sync.RWMutex{},
		buckets: map[uint64]*bucket[StringKey, string]{},
	}

	res, ok := shard.get("doesnotexist")

	assert.Equal(t, ok, false)
	assert.Equal(t, res, "")
}

func TestDelete(t *testing.T) {
	shard := &shard[StringKey, string]{
		m:       sync.RWMutex{},
		buckets: map[uint64]*bucket[StringKey, string]{},
	}

	shard.set("abc", "def")
	res, ok := shard.get("abc")
	assert.Equal(t, ok, true)
	assert.Equal(t, res, "def")

	shard.delete("abc")
	res, ok = shard.get("abc")

	assert.Equal(t, ok, false)

	// We expect the zero-value of the key type to be returned
	assert.Equal(t, res, "")
}
