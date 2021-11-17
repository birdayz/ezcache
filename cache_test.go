package ezcache

import (
	"errors"
	"testing"

	"gotest.tools/v3/assert"
)

type KeyString string

func (ks KeyString) Equals(s KeyString) bool {
	return s == ks
}

func (ks KeyString) HashCode() uint64 {
	return 0
}

func TestCacheLoaderLoads(t *testing.T) {
	cache := New(func(key KeyString) (string, error) { return "test-val", nil }, 10)

	res, err := cache.Get("my-key")
	assert.NilError(t, err)
	assert.Equal(t, res, "test-val")
}

func TestCacheLoaderError(t *testing.T) {
	loaderError := "could not connect to database"
	cache := New(func(key KeyString) (string, error) { return "", errors.New(loaderError) }, 10)

	_, err := cache.Get("my-key")
	assert.ErrorContains(t, err, loaderError)
}

func TestCacheDelete(t *testing.T) {

	cache := New(func(key KeyString) (string, error) { return "", errors.New("dont want to load") }, 10)
	cache.Set("key", "value")
	res, err := cache.Get("key")
	assert.NilError(t, err)
	assert.Equal(t, res, "value")
	cache.Delete("key")
	res, err = cache.Get("key")
	assert.ErrorContains(t, err, "dont want to")
}
