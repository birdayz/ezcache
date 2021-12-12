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


## TBD
- Cache errors, when would you want to have this?
- LoaderFn where previous value is known
- What about ctx?
- Reduce number of hash func calls
- Allow unlimited capacity?
- Allow unlimited TTL
