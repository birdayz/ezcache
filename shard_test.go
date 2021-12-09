package ezcache

import (
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestGet(t *testing.T) {
	shard := newShard[StringKey, string](10)
	abc := StringKey("abc")

	shard.set("abc", abc.HashCode(), "def")
	res, ok := shard.get("abc", abc.HashCode())

	assert.Equal(t, ok, true)
	assert.Equal(t, res, "def")
}

func TestSetGetEvict(t *testing.T) {
	shard := newShard[StringKey, string](2)
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
	shard := newShard[StringKey, string](2)

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
	shard := newShard[StringKey, string](10)

	doesnotexist := StringKey("doesnotexist")
	res, ok := shard.get("doesnotexist", doesnotexist.HashCode())

	assert.Equal(t, ok, false)
	assert.Equal(t, res, "")
}

func TestDelete(t *testing.T) {
	shard := newShard[StringKey, string](10)

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

type fake struct {
	key      string
	hashCode uint64
}

func (f fake) Equals(f2 fake) bool {
	return f.key == f2.key
}

func (f fake) HashCode() uint64 { return f.hashCode }

func TestDeleteSameHashCode(t *testing.T) {
	shard := newShard[fake, string](10)

	abc := fake{"abc", 0}
	shard.set(abc, abc.HashCode(), "val1")
	res, ok := shard.get(abc, abc.HashCode())
	assert.Equal(t, ok, true)
	assert.Equal(t, res, "val1")

	def := fake{"def", 0}
	shard.set(def, def.HashCode(), "val2")
	res, ok = shard.get(def, def.HashCode())
	assert.Equal(t, ok, true)
	assert.Equal(t, res, "val2")

	// val1 should be still working after a value on same bucket was inserted
	res, ok = shard.get(abc, abc.HashCode())
	assert.Equal(t, ok, true)
	assert.Equal(t, res, "val1")

	// Delete abc, it should now be gone
	shard.delete(abc)
	res, ok = shard.get(abc, abc.HashCode())
	assert.Equal(t, ok, false)

	// Deleting def works as well
	shard.delete(def)
	res, ok = shard.get(def, def.HashCode())
	assert.Equal(t, ok, false)
}

func TestExpireTTL(t *testing.T) {
	shard := newShard[StringKey, string](10)
	shard.ttl = time.Millisecond * 10

	var fakeTime time.Time

	// "Inject" fake time
	timeFn := func() time.Time {
		return fakeTime
	}
	timeNow = timeFn
	defer func() { timeNow = time.Now }()

	fakeTime = time.Now()

	abc := StringKey("abc")
	shard.set("abc", abc.HashCode(), "def")

	fakeTime = fakeTime.Add(time.Millisecond * 11)

	_, ok := shard.get(abc, abc.HashCode())
	assert.Equal(t, ok, false)

}

func TestExpireTTLProlongedAfterSet(t *testing.T) {
	shard := newShard[StringKey, string](10)
	shard.ttl = time.Millisecond * 10

	abc := StringKey("abc")
	def := StringKey("def")
	shard.set("abc", abc.HashCode(), "def")
	shard.set(def, def.HashCode(), "def2")

	time.Sleep(time.Millisecond * 5) // TODO, replace timing based tests with mocked/mock-able clock
	_, ok := shard.get("abc", abc.HashCode())
	assert.Equal(t, ok, true)

	shard.set("abc", abc.HashCode(), "defNew")
	time.Sleep(time.Millisecond * 5) // TODO, replace timing based tests with mocked/mock-able clock

	// Check if the first item, which was touched, is still around
	_, ok = shard.get("abc", abc.HashCode())
	assert.Equal(t, ok, true)

	// Check if the second item, which we did not touch, was remove successfully
	_, ok = shard.get(def, def.HashCode())
	assert.Equal(t, ok, false)
}
