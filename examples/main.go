package main

import (
	"errors"
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
	a := ezcache.New(
		func(key *TestKey) (*[]string, error) {
			x := []string{"value"}
			return &x, nil
		},
		2,
	)

	k := &TestKey{
		blah: 0,
		bleh: "",
	}

	a.Set(k, &[]string{"my-value"})

	res, err := a.Get(k)
	fmt.Println(res, err)

	{
		b := ezcache.New(func(key ezcache.StringKey) (string, error) { return "dsad", errors.New("fail") }, 256)

		back, _ := b.Get("abc")
		fmt.Println("Got back", back)
	}
}
