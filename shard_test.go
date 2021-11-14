package ezcache

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestGet(t *testing.T) {
	shard := &shard[string, string]{
		buckets: map[uint64]bucket[string, string]{},
		hasher:  StringHasher,
	}

	shard.set("abc", "def")
	res := shard.get("abc")

	assert.Equal(t, res, "def")
}

func TestGetDoesNotExist(t *testing.T) {
	shard := &shard[string, string]{
		buckets: map[uint64]bucket[string, string]{},
		hasher:  StringHasher,
	}

	res := shard.get("doesnotexist")

	assert.Equal(t, res, "")
}

func TestDelete(t *testing.T) {
	shard := &shard[string, string]{
		buckets: map[uint64]bucket[string, string]{},
		hasher:  StringHasher,
	}

	shard.set("abc", "def")
	res := shard.get("abc")
	assert.Equal(t, res, "def")

	shard.delete("abc")
	res = shard.get("abc")

	// We expect the zero-value of the key type to be returned
	assert.Equal(t, res, "")
}
