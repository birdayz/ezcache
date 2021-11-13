package main

import (
	"fmt"
	"hash/fnv"

	"github.com/birdayz/ezcache"
)

type TestKey struct {
	blah int
	bleh string
	x    map[string]string
}

func (t TestKey) String() string {
	return t.bleh

}

func main() {
	r := fmt.Sprintf("%v", TestKey{blah: 1, bleh: "x"})
	fmt.Println(r)
	a := ezcache.New(
		func(key *TestKey) (string, error) {
			return "value", nil
		},
		func(key *TestKey) uint32 {
			h := fnv.New32a()
			h.Write([]byte(key.bleh))
			return h.Sum32()
		},
	)

	k := &TestKey{
		blah: 0,
		bleh: "",
	}

	a.Set(k, "my-value")

	res, err := a.Get(k)
	fmt.Println(res, err)

}
