//+build metago

// function with struct processing.

package simpleFunc

import (
	"fmt"

	"github.com/vvakame/til/go/metago"
)

type Foo struct {
	ID   int64
	Name string
}

func example() {
	{
		obj := &Foo{}
		kvMap := make(map[string]interface{})
		mv := metago.ValueOf(obj)
		for _, mf := range mv.Fields() {
			kvMap[mf.Name()] = mf.Value()
		}
		fmt.Println(kvMap)
	}
	{
		obj := Foo{}
		kvMap := make(map[string]interface{})
		mv := metago.ValueOf(obj)
		for _, mf := range mv.Fields() {
			kvMap[mf.Name()] = mf.Value()
		}
		fmt.Println(kvMap)
	}
	{
		var obj *Foo
		kvMap := make(map[string]interface{})
		mv := metago.ValueOf(obj)
		for _, mf := range mv.Fields() {
			kvMap[mf.Name()] = mf.Value()
		}
		fmt.Println(kvMap)
	}
	{
		var obj Foo
		kvMap := make(map[string]interface{})
		mv := metago.ValueOf(obj)
		for _, mf := range mv.Fields() {
			kvMap[mf.Name()] = mf.Value()
		}
		fmt.Println(kvMap)
	}
	{
		obj := new(Foo)
		kvMap := make(map[string]interface{})
		mv := metago.ValueOf(obj)
		for _, mf := range mv.Fields() {
			kvMap[mf.Name()] = mf.Value()
		}
		fmt.Println(kvMap)
	}
}
