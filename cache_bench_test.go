package ezcache

import "testing"

func BenchmarkSet(b *testing.B) {

	cache := New[StringKey, StringKey](nil, 1, 1000)
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			cache.Set("abc", "val")
		}
	})

}
