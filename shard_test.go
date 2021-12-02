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

	abc := StringKey("abc")

	shard.set("abc", abc.HashCode(), "def")
	res, ok := shard.get("abc", abc.HashCode())

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

	first := StringKey("first")

	shard.set("first", first.HashCode(), "def")
	res, ok := shard.get("first", first.HashCode())
	assert.Equal(t, ok, true)
	assert.Equal(t, res, "def")

	second := StringKey("second")
	shard.set(second, second.HashCode(), "second_val")
	res, ok = shard.get("second", second.HashCode())
	assert.Equal(t, ok, true)
	assert.Equal(t, res, "second_val")

	// Trigger eviction
	third := StringKey("third")
	shard.set(third, third.HashCode(), "third_val")
	res, ok = shard.get("third", third.HashCode())
	assert.Equal(t, ok, true)
	assert.Equal(t, res, "third_val", third.HashCode())

	res, ok = shard.get("first", first.HashCode())
	assert.Equal(t, ok, false)
}

func TestSetGetEvictOrder(t *testing.T) {
	shard := &shard[StringKey, string]{
		m:          sync.RWMutex{},
		buckets:    map[uint64]*bucket[StringKey, string]{},
		linkedList: NewList[StringKey](),
		capacity:   2,
	}

	first := StringKey("first")
	shard.set("first", first.HashCode(), "def")
	res, ok := shard.get("first", first.HashCode())
	assert.Equal(t, ok, true)
	assert.Equal(t, res, "def")

	second := StringKey("second")
	shard.set("second", second.HashCode(), "second_val")
	res, ok = shard.get("second", second.HashCode())
	assert.Equal(t, ok, true)
	assert.Equal(t, res, "second_val")

	// Issue get, so first is recently used; shall not be evicted
	shard.get("first", first.HashCode())

	// Trigger eviction
	third := StringKey("third")
	shard.set("third", third.HashCode(), "third_val")
	res, ok = shard.get("third", third.HashCode())
	assert.Equal(t, ok, true)
	assert.Equal(t, res, "third_val")

	// First has to be preserved, because we touched it
	res, ok = shard.get("first", first.HashCode())
	assert.Equal(t, ok, true)

	res, ok = shard.get("second", second.HashCode())
	assert.Equal(t, ok, false)
}

func TestGetDoesNotExist(t *testing.T) {
	shard := &shard[StringKey, string]{
		m:          sync.RWMutex{},
		buckets:    map[uint64]*bucket[StringKey, string]{},
		linkedList: NewList[StringKey](),
		capacity:   10,
	}

	doesnotexist := StringKey("doesnotexist")
	res, ok := shard.get("doesnotexist", doesnotexist.HashCode())

	assert.Equal(t, ok, false)
	assert.Equal(t, res, "")
}

func TestDelete(t *testing.T) {
	shard := &shard[StringKey, string]{
		m:          sync.RWMutex{},
		buckets:    map[uint64]*bucket[StringKey, string]{},
		linkedList: NewList[StringKey](),
		capacity:   10,
	}

	abc := StringKey("abc")
	shard.set("abc", abc.HashCode(), "def")
	res, ok := shard.get("abc", abc.HashCode())
	assert.Equal(t, ok, true)
	assert.Equal(t, res, "def")

	shard.delete("abc")
	res, ok = shard.get("abc", abc.HashCode())

	assert.Equal(t, ok, false)

	// We expect the zero-value of the key type to be returned
	assert.Equal(t, res, "")
}
