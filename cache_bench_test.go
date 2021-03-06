package ezcache

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"testing"
)

func BenchmarkSetString(b *testing.B) {
	cfg := NewBuilder[StringKey, string]().Capacity(b.N / 2).NumShards(1)

	b.ResetTimer()
	b.Run("Set", func(b *testing.B) {
		var oldMemProfileRate = runtime.MemProfileRate
		runtime.MemProfileRate = 0

		cache := cfg.Build()

		b.ResetTimer()
		runtime.MemProfileRate = oldMemProfileRate

		for i := 0; i < b.N; i++ {
			cache.Set(StringKey(strconv.Itoa(i)), "")
		}
	})
	b.Run("Get", func(b *testing.B) {
		var oldMemProfileRate = runtime.MemProfileRate
		runtime.MemProfileRate = 0

		cache := cfg.Build()

		for i := 0; i < b.N; i++ {
			cache.Set(StringKey(strconv.Itoa(i)), strconv.Itoa(i))
		}

		b.ResetTimer()
		runtime.MemProfileRate = oldMemProfileRate

		for i := 0; i < b.N; i++ {
			_, _ = cache.Get(StringKey(strconv.Itoa(i)))
		}
	})
}

func BenchmarkSetStringNotFound(b *testing.B) {
	cache := NewBuilder[StringKey, string]().Capacity(b.N).NumShards(1).Build()

	for i := 0; i < b.N; i++ {
		_, _ = cache.Get(StringKey(strconv.Itoa(i)))
	}

}

func BenchmarkSetInt(b *testing.B) {
	cfg := NewBuilder[IntKey, int]().Capacity(10).NumShards(1)

	b.ResetTimer()
	b.Run("Set", func(b *testing.B) {
		var oldMemProfileRate = runtime.MemProfileRate
		runtime.MemProfileRate = 0

		cache := cfg.Build()

		b.ResetTimer()
		runtime.MemProfileRate = oldMemProfileRate

		for i := 0; i < b.N; i++ {
			cache.Set(IntKey(i), i)
		}
	})
	b.Run("Get", func(b *testing.B) {
		var oldMemProfileRate = runtime.MemProfileRate
		runtime.MemProfileRate = 0

		cache := cfg.Build()

		for i := 0; i < b.N; i++ {
			cache.Set(IntKey(i), i)
		}

		b.ResetTimer()
		runtime.MemProfileRate = oldMemProfileRate

		for i := 0; i < b.N; i++ {
			_, _ = cache.Get(IntKey(i))
		}
	})
}

func BenchmarkParallelSet(b *testing.B) {
	tests := []struct {
		parallelism    int
		itemsPerWorker int
		buckets        int
	}{
		{
			parallelism:    1,
			itemsPerWorker: 100000,
			buckets:        1,
		},
		{
			parallelism:    2,
			itemsPerWorker: 50000,
			buckets:        1,
		},
		{
			parallelism:    10,
			itemsPerWorker: 10000,
			buckets:        1,
		},
		{
			parallelism:    100,
			itemsPerWorker: 1000,
			buckets:        1,
		},
		{
			parallelism:    1,
			itemsPerWorker: 100000,
			buckets:        256,
		},
		{
			parallelism:    2,
			itemsPerWorker: 50000,
			buckets:        256,
		},
		{
			parallelism:    10,
			itemsPerWorker: 10000,
			buckets:        256,
		},
		{
			parallelism:    100,
			itemsPerWorker: 1000,
			buckets:        256,
		},
		{
			parallelism:    100,
			itemsPerWorker: 100000,
			buckets:        256,
		},
	}

	for _, tt := range tests {

		b.Run(fmt.Sprintf("%v-Parallel-%v-Buckets-%v", tt.itemsPerWorker*tt.parallelism, tt.parallelism, tt.buckets), func(b *testing.B) {
			parallelism := tt.parallelism
			itemsPerWorker := tt.itemsPerWorker
			buckets := tt.buckets

			cache := NewBuilder[IntKey, StringKey]().NumShards(buckets).Capacity(100000).Build()

			data := make(map[int][]int) // one entry per worker

			for i := 0; i < parallelism; i++ {
				data[i] = make([]int, 0, itemsPerWorker)

				for d := 0; d < itemsPerWorker; d++ {
					data[i] = append(data[i], rand.Int())
				}
			}

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				// Do actual test
				var wg sync.WaitGroup

				for w := 0; w < parallelism; w++ {
					wg.Add(1)
					go func(workerID int) {

						workerItems := data[workerID]
						for _, workerItem := range workerItems {
							cache.Set(IntKey(workerItem), "value")
						}

						wg.Done()
					}(w)
				}

				wg.Wait()

			}

		})
	}
}
