package ezcache

import (
	"testing"
)

func BenchmarkHashMapSet(b *testing.B) {

	b.Run("Set", func(b *testing.B) {
		m := NewHashMap[IntKey, int](16)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			m.Set(IntKey(i), i)
		}
	})
	b.Run("Get", func(b *testing.B) {
		m := NewHashMap[IntKey, int](16)
		for i := 0; i < b.N; i++ {
			m.Set(IntKey(i), i)
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			res, _ := m.Get(IntKey(i))
			if res != i {
				b.FailNow()
			}

		}
	})
}

func BenchmarkGoMapSet(b *testing.B) {

	b.Run("Set", func(b *testing.B) {
		m := make(map[int]int)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			m[i] = i
		}
	})
	b.Run("Get", func(b *testing.B) {
		m := make(map[int]int)
		for i := 0; i < b.N; i++ {
			m[i] = i
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			res := m[i]
			if res != i {
				b.FailNow()
			}
		}
	})
}
