package ezcache

import (
	"errors"
	"testing"

	"gotest.tools/v3/assert"
)

func TestCacheLoaderLoads(t *testing.T) {
	cache := New(func(key string) (string, error) { return "test-val", nil }, StringHasher, 10)

	res, err := cache.Get("my-key")
	assert.NilError(t, err)
	assert.Equal(t, res, "test-val")
}

func TestCacheLoaderError(t *testing.T) {
	loaderError := "could not connect to database"
	cache := New(func(key string) (string, error) { return "", errors.New(loaderError) }, StringHasher, 10)

	_, err := cache.Get("my-key")
	assert.ErrorContains(t, err, loaderError)
}
