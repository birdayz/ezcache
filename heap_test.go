package ezcache

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestHeap(t *testing.T) {
	h := NewHeap(AscendingComparator[int])

	h.Push(100)
	h.Push(1)
	h.Push(5)
	h.Push(2)

	peekRes := h.Peek()
	assert.Equal(t, peekRes.Item, 1)
	res := h.Pop()
	assert.Equal(t, res.Item, 1)

	peekRes = h.Peek()
	assert.Equal(t, peekRes.Item, 2)

	res = h.Pop()
	assert.Equal(t, res.Item, 2)

	peekRes = h.Peek()
	assert.Equal(t, peekRes.Item, 5)
	res = h.Pop()
	assert.Equal(t, res.Item, 5)

	peekRes = h.Peek()
	assert.Equal(t, peekRes.Item, 100)
	res = h.Pop()
	assert.Equal(t, res.Item, 100)

}

func TestHeapEmpty(t *testing.T) {
	t.Skip()
	h := NewHeap(AscendingComparator[int])

	h.Peek()
	h.Pop()
}
