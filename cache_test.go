package ezcache

import (
	"errors"
	"testing"

	"gotest.tools/v3/assert"
)

func TestCacheLoaderLoads(t *testing.T) {
	cache := New(func(key StringKey) (string, error) { return "test-val", nil }, 10, 10)

	res, err := cache.Get("my-key")
	assert.NilError(t, err)
	assert.Equal(t, res, "test-val")
}

func TestCacheLoaderError(t *testing.T) {
	loaderError := "could not connect to database"
	cache := New(func(key StringKey) (string, error) { return "", errors.New(loaderError) }, 10, 10)

	_, err := cache.Get("my-key")
	assert.ErrorContains(t, err, loaderError)
}

func TestCacheDelete(t *testing.T) {
	cache := New(func(key StringKey) (string, error) { return "", errors.New("dont want to load") }, 10, 10)
	cache.Set("key", "value")
	res, err := cache.Get("key")
	assert.NilError(t, err)
	assert.Equal(t, res, "value")
	cache.Delete("key")
	res, err = cache.Get("key")
	assert.ErrorContains(t, err, "dont want to")
}

func TestCacheSetSet(t *testing.T) {
	cache := New(func(key StringKey) (string, error) { return "", errors.New("dont want to load") }, 10, 10)
	cache.Set("key", "value")
	cache.Set("key", "value2")
	res, err := cache.Get("key")
	assert.NilError(t, err)
	assert.Equal(t, res, "value2")

}

func TestExpire(t *testing.T) {
	cache := New[StringKey, StringKey](nil, 10, 2)

	cache.Set("key1", "val1")
	cache.Set("key2", "val2")
	cache.Set("key3", "val3")

	res, err := cache.Get("key1")
	assert.Error(t, err, ErrNotFound.Error())
	assert.Equal(t, res, StringKey(""))
}

func TestExpireOrder(t *testing.T) {
	cache := New[StringKey, StringKey](nil, 10, 2)

	cache.Set("key1", "val1")
	cache.Set("key2", "val2")
	cache.Set("key3", "val3")
	cache.Set("key4", "val4")

	_, err := cache.Get("key1")
	assert.Error(t, err, ErrNotFound.Error())

	_, err = cache.Get("key2")
	assert.Error(t, err, ErrNotFound.Error())

	res, err := cache.Get("key3")
	assert.NilError(t, err)
	assert.Equal(t, res, StringKey("val3"))

	res, err = cache.Get("key4")
	assert.NilError(t, err)
	assert.Equal(t, res, StringKey("val4"))
}
