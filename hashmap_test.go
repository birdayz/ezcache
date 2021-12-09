package ezcache

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestHashMapSet(t *testing.T) {

	m := NewHashMap[StringKey, string]()

	m.Set("keya", "valuea")

	res, ok := m.Get("keya")
	assert.Equal(t, ok, true)
	assert.Equal(t, res, "valuea")

	setRes := m.Set("keya", "valueb")
	assert.Equal(t, setRes, true)

	res, ok = m.Get("keya")
	assert.Equal(t, ok, true)
	assert.Equal(t, res, "valueb")
}
