package ezcache

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

func BenchmarkSetString(b *testing.B) {

	b.ResetTimer()
	b.Run("Set", func(b *testing.B) {
		cache := New[StringKey, string](nil, 100, 10000)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			cache.Set(StringKey(strconv.Itoa(i)), "")
		}
	})
	b.Run("Get", func(b *testing.B) {
		cache := New[StringKey, string](nil, 100, 10000)

		for i := 0; i < b.N; i++ {
			cache.Set(StringKey(strconv.Itoa(i)), strconv.Itoa(i))
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = cache.Get(StringKey(strconv.Itoa(i)))
		}
	})
}

func BenchmarkSetInt(b *testing.B) {

	b.ResetTimer()
	b.Run("Set", func(b *testing.B) {
		cache := New[IntKey, int](nil, 10, 16)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			cache.Set(IntKey(i), i)
		}
	})
	b.Run("Get", func(b *testing.B) {
		cache := New[IntKey, int](nil, 1, 10000)
		for i := 0; i < b.N; i++ {
			cache.Set(IntKey(i), i)
		}

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			cache.Get(IntKey(i))
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

			cache := New[IntKey, StringKey](nil, buckets, 100000)

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
