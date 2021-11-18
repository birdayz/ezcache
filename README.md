## Goals
- [x] Cache Loader support 
- Capacity
- TTL
- [ ] Write through
- Cache Loader with current entry (-> Multiple implementations for the cache loader possible)
- Bulk Loader
- TBD: Sync vs Async load. What would it mean?
- Refresh Ahead
- Invalidation via stream?
- Cache errors?
- Metrics
- Capacity
- Work with binary data (for both in-mem and semi persistent storage); employ marshalers/unmarshalers
- Consider using mmap for the storage
- Security,Safety: There are no collisions, i.e. if you Set/Get a key, you will never get a different key's value back because (with low chance) some other key got the same hash. Not even with low odds - not at all. We consider this extremely important for multi-tenant systems, where even rare collisions can comprimise other user's data.
- No / only very "standard" and vetted deps (i.e. pkg/errors, go/x)
- Sub-package for practical, ready to use modules that use the cache, i.e.: grpc caching middleware, etag based cache for http, ...


## TBD
- Cache errors, when would you want to have this?
- LoaderFn where previous value is known
- What about ctx?
