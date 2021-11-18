package ezcache

import "hash/fnv"

type StringKey string

func (ks StringKey) Equals(s StringKey) bool {
	return s == ks
}

func (ks StringKey) HashCode() uint64 {
	h := fnv.New64a()
	h.Write([]byte(ks))
	return h.Sum64()
}
