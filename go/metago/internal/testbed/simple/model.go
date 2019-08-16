//+build metago

package main

import (
	"fmt"

	"github.com/vvakame/til/go/metago"
)

type Foo struct {
	ID   int64
	Name string
}

func main() {
	obj := &Foo{1, "vvakame"}
	mv := metago.ValueOf(obj)
	for _, mf := range mv.Fields() {
		fmt.Println(mf.Name(), mf.Value())
	}
}
