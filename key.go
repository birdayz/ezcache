package ezcache

import (
	"hash/fnv"
)

type StringKey string

func (ks StringKey) Equals(s StringKey) bool {
	return s == ks
}

func (ks StringKey) HashCode() uint64 {
	h := fnv.New64a()
	h.Write([]byte(ks))
	return h.Sum64()
}

type IntKey int

func (ik IntKey) Equals(s IntKey) bool {
	return ik == s
}

func (ik IntKey) HashCode() uint64 {
	return uint64(ik)
}
