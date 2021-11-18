package ezcache

import (
	"errors"
	"testing"

	"gotest.tools/v3/assert"
)

func TestCacheLoaderLoads(t *testing.T) {
	cache := New(func(key StringKey) (string, error) { return "test-val", nil }, 10)

	res, err := cache.Get("my-key")
	assert.NilError(t, err)
	assert.Equal(t, res, "test-val")
}

func TestCacheLoaderError(t *testing.T) {
	loaderError := "could not connect to database"
	cache := New(func(key StringKey) (string, error) { return "", errors.New(loaderError) }, 10)

	_, err := cache.Get("my-key")
	assert.ErrorContains(t, err, loaderError)
}

func TestCacheDelete(t *testing.T) {
	cache := New(func(key StringKey) (string, error) { return "", errors.New("dont want to load") }, 10)
	cache.Set("key", "value")
	res, err := cache.Get("key")
	assert.NilError(t, err)
	assert.Equal(t, res, "value")
	cache.Delete("key")
	res, err = cache.Get("key")
	assert.ErrorContains(t, err, "dont want to")
}

func TestCacheSetSet(t *testing.T) {
	cache := New(func(key StringKey) (string, error) { return "", errors.New("dont want to load") }, 10)
	cache.Set("key", "value")
	cache.Set("key", "value2")
	res, err := cache.Get("key")
	assert.NilError(t, err)
	assert.Equal(t, res, "value2")

}
