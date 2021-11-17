## Goals
- Cache Loader support
- Cache Loader with current entry (-> Multiple implementations for the cache loader possible)
- Bulk Loader
- TBD: Sync vs Async load. What would it mean?
- Refresh Ahead
- Invalidation via stream?
- Cache errors?
- Metrics
- Capacity
- TTL
- Work with binary data (for both in-mem and semi persistent storage); employ marshalers/unmarshalers
- Consider using mmap for the storage


## TBD
- Cache errors, when would you want to have this?
- LoaderFn where previous value is known
- What about ctx?
