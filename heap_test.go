package ezcache

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestHeap(t *testing.T) {
	h := NewHeap(ComparableComparator[int])

	h.Push(100)
	h.Push(1)
	h.Push(5)
	h.Push(2)

	res := h.Pop()
	assert.Equal(t, res, 1)

	res = h.Pop()
	assert.Equal(t, res, 2)

	res = h.Pop()
	assert.Equal(t, res, 5)

	res = h.Pop()
	assert.Equal(t, res, 100)
}
