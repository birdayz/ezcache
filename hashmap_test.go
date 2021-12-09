package ezcache

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestHashMapSet(t *testing.T) {
	m := NewHashMap[StringKey, string](16)

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

func TestHashMapSetWithGrow(t *testing.T) {
	m := NewHashMap[IntKey, int](16)

	for i := 0; i < 5000; i++ {
		m.Set(IntKey(i), i)

		res, ok := m.Get(IntKey(i))
		assert.Equal(t, ok, true)
		assert.Equal(t, res, i)
	}

}
