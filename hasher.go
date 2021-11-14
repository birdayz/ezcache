package ezcache

import (
	"fmt"
	"hash/fnv"
)

var StringHasher = func(key string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(key))
	return h.Sum64()

}

var DefaultHasher = func(key string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%v", key)))
	return h.Sum64()
}
