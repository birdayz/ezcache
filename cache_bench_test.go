package ezcache

import "testing"

func BenchmarkSet(b *testing.B) {
	cache := New[StringKey, StringKey](nil, 1, 1000)
	for i := 0; i < b.N; i++ {
		cache.Set("abc", "val")
	}

}
