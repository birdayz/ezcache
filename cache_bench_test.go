package ezcache

import (
	"fmt"
	"sync"
	"testing"

	"github.com/thanhpk/randstr"
)

func BenchmarkSet(b *testing.B) {
	cache := New[StringKey, StringKey](nil, 256, 1000)
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			cache.Set("abc", "val")
		}
	})

}

func BenchmarkParallelSet(b *testing.B) {
	table := []struct {
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
	}

	for _, tt := range table {

		b.Run(fmt.Sprintf("Parallel-%v-Buckets-%v", tt.parallelism, tt.buckets), func(b *testing.B) {
			parallelism := tt.parallelism
			itemsPerWorker := tt.itemsPerWorker
			buckets := tt.buckets

			cache := New[StringKey, StringKey](nil, buckets, 1000)

			data := make(map[int][]string) // one entry per worker

			for i := 0; i < parallelism; i++ {
				data[i] = make([]string, 0, itemsPerWorker)

				for d := 0; d < itemsPerWorker; d++ {
					data[i] = append(data[i], randstr.String(10))
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
							cache.Set(StringKey(workerItem), "value")
						}

						wg.Done()
					}(w)
				}

				wg.Wait()

			}

		})
	}
}
