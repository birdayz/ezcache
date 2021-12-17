## Goals
- [x] Cache Loader support 
- [x] Capacity
- [x] TTL
- [ ] Write through
- Cache Loader with current entry (-> Multiple implementations for the cache loader possible)
- Bulk Loader
- TBD: Sync vs Async load. What would it mean?
- Refresh Ahead
- Invalidation via stream?
- Cache errors?
- Metrics
- Capacity
- Security,Safety: There are no collisions, i.e. if you Set/Get a key, you will never get a different key's value back because (with low chance) some other key got the same hash. Not even with low odds - not at all. We consider this extremely important for multi-tenant systems, where even rare collisions can comprimise other user's data.
- Sub-package for practical, ready to use modules that use the cache, i.e.: grpc caching middleware, etag based cache for http, ...
- Eviction callbacks
- Admission policy? (See TinyLFU paper)
- Generic map mit equals + hashcode !
- Generic singleflight
- Equals() + HashCode() vs comparable + hashCode -> use equals ?!
- Refactor shard.go to not use map, but use slice/array + linkedlist
- Frequently used vs Recently; auto refresh top frequently used? frequently more realistic?
- refresh only frequent, but not recent data?
- bloomfilter
- copy key?
- More flexible expiry + evict + admission policy
- Split expiration into specific type, to only use it if needed. Expiration manager?
- CacheLoader: reload
- Singleflight
- Expiration optional


## TBD/TODO
- Cache errors, when would you want to have this?
- LoaderFn where previous value is known
- What about ctx?
- Reduce number of hash func calls
- Allow unlimited capacity?
- Allow unlimited TTL
- Try to not alloc for LinkedList/Heap - we could use the same pointer as for bucketItem, but add relevant methods for LL/Heap operations.
- LinkedList should be pointer to bucketItem ? then deletes can avoid hashing key again.
- Cached timer! every 1s/ms

## Performance

```
go test -benchtime=10000000x -run='^$' -bench=BenchmarkSet -benchmem -memprofile memprofile.out -cpuprofile profile.out -count=5
goos: linux
goarch: amd64
pkg: github.com/birdayz/ezcache
cpu: AMD Ryzen 9 3900X 12-Core Processor
BenchmarkSetString/Set-24         	10000000	       691.6 ns/op	     276 B/op	       7 allocs/op
BenchmarkSetString/Set-24         	10000000	       684.9 ns/op	     276 B/op	       7 allocs/op
BenchmarkSetString/Set-24         	10000000	       701.2 ns/op	     276 B/op	       7 allocs/op
BenchmarkSetString/Set-24         	10000000	       711.4 ns/op	     276 B/op	       7 allocs/op
BenchmarkSetString/Set-24         	10000000	       701.0 ns/op	     276 B/op	       7 allocs/op
BenchmarkSetString/Get-24         	10000000	       259.6 ns/op	      23 B/op	       1 allocs/op
BenchmarkSetString/Get-24         	10000000	       253.7 ns/op	      23 B/op	       1 allocs/op
BenchmarkSetString/Get-24         	10000000	       260.4 ns/op	      23 B/op	       1 allocs/op
BenchmarkSetString/Get-24         	10000000	       251.7 ns/op	      23 B/op	       1 allocs/op
BenchmarkSetString/Get-24         	10000000	       257.0 ns/op	      23 B/op	       1 allocs/op
BenchmarkSetInt/Set-24            	10000000	       439.6 ns/op	     208 B/op	       6 allocs/op
BenchmarkSetInt/Set-24            	10000000	       441.9 ns/op	     208 B/op	       6 allocs/op
BenchmarkSetInt/Set-24            	10000000	       445.4 ns/op	     208 B/op	       6 allocs/op
BenchmarkSetInt/Set-24            	10000000	       431.5 ns/op	     208 B/op	       6 allocs/op
BenchmarkSetInt/Set-24            	10000000	       439.8 ns/op	     208 B/op	       6 allocs/op
BenchmarkSetInt/Get-24            	10000000	        82.50 ns/op	       7 B/op	       0 allocs/op
BenchmarkSetInt/Get-24            	10000000	        79.30 ns/op	       7 B/op	       0 allocs/op
BenchmarkSetInt/Get-24            	10000000	        79.43 ns/op	       8 B/op	       0 allocs/op
BenchmarkSetInt/Get-24            	10000000	        82.02 ns/op	       7 B/op	       0 allocs/op
BenchmarkSetInt/Get-24            	10000000	        82.85 ns/op	       7 B/op	       0 allocs/op
PASS
ok  	github.com/birdayz/ezcache	134.206s
```


