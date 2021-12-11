package main

import (
	"fmt"
	"reflect"

	"github.com/birdayz/ezcache"
)

type TestKey struct {
	blah int
	bleh string
}

func (t TestKey) String() string {
	return t.bleh
}

func (t *TestKey) Equals(other *TestKey) bool {
	return reflect.DeepEqual(t, other)
}

func (t *TestKey) HashCode() uint64 {
	return uint64(0)
}

func main() {
	cache := ezcache.NewBuilder[*TestKey, []string]().Capacity(10).NumShards(2).Build()

	k := &TestKey{
		blah: 0,
		bleh: "",
	}

	cache.Set(k, []string{"my-value"})

	res, err := cache.Get(k)
	fmt.Println(res, err)
}
