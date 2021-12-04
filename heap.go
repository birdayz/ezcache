package ezcache

import "constraints"

type HeapElement[T any] struct {
	index int
	Item  T
}

type Heap[T any] struct {
	comparator Comparator[T]
	data       []*HeapElement[T]
}

func NewHeap[T any](comparator Comparator[T]) *Heap[T] {
	return &Heap[T]{
		comparator: comparator,
		data:       make([]*HeapElement[T], 0, 0),
	}
}

func AscendingComparator[T constraints.Ordered](t1, t2 T) int {
	if t1 > t2 {
		return 1
	} else if t1 < t2 {
		return -1
	} else {
		return 0
	}
}

type Comparator[T any] func(t1, t2 T) int

func (t *Heap[T]) Peek() (item *HeapElement[T]) {
	return t.data[0]
}

func (t *Heap[T]) Init() {
	n := len(t.data)
	for i := n/2 - 1; i >= 0; i-- {
		t.down(i, n)
	}

}

func (t *Heap[T]) Fix(he *HeapElement[T]) {
	i := he.index
	if i == -1 {
		return
	}
	if !t.down(i, len(t.data)) {
		t.up(i)
	}
}

func (t *Heap[T]) Push(item T) *HeapElement[T] {
	hi := &HeapElement[T]{
		Item: item,
	}
	t.data = append(t.data, hi)
	t.up(len(t.data) - 1)

	hi.index = len(t.data) - 1

	return hi
}

func (t *Heap[T]) Pop() (item *HeapElement[T]) {
	{
		n := len(t.data) - 1

		t.data[0], t.data[n] = t.data[n], t.data[0]
		t.data[0].index = 0
		t.data[n].index = n

		t.down(0, n)
	}

	old := t.data
	n := len(old)
	x := old[n-1]
	t.data = old[0 : n-1]
	x.index = -1
	return x
}

func (t *Heap[T]) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || t.comparator(t.data[j].Item, t.data[i].Item) > 0 {
			break
		}
		t.data[i], t.data[j] = t.data[j], t.data[i]
		t.data[i].index = i
		t.data[j].index = j

		j = i
	}
}

func (t *Heap[T]) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && t.comparator(t.data[j2].Item, t.data[j1].Item) < 0 {
			j = j2 // = 2*i + 2  // right child
		}
		if t.comparator(t.data[j].Item, t.data[i].Item) > 0 {
			break
		}

		t.data[i], t.data[j] = t.data[j], t.data[i]
		t.data[i].index = i
		t.data[j].index = j

		i = j
	}
	return i > i0
}

func (t *Heap[T]) Remove(he *HeapElement[T]) {
	i := he.index
	if i == -1 {
		return
	}

	{
		n := len(t.data) - 1
		if n != i {

			t.data[i], t.data[n] = t.data[n], t.data[i]
			t.data[i].index = i
			t.data[n].index = n

			if !t.down(i, n) {
				t.up(i)
			}
		}
	}

	// Pop
	old := t.data
	n := len(old)
	x := old[n-1]
	t.data = old[0 : n-1]
	x.index = -1
}
