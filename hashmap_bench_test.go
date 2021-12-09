package ezcache

import (
	"strconv"
	"testing"
)

func BenchmarkHashMapSetString(b *testing.B) {
	m := NewHashMap[StringKey, StringKey]()

	b.ResetTimer()

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			m.Set(StringKey(strconv.Itoa(i)), "")
		}
	})
	// b.Run("Get", func(b *testing.B) {
	// 	for i := 0; i < b.N; i++ {
	// 		cache.Get(StringKey(strconv.Itoa(i)))
	// 	}
	// })
}
