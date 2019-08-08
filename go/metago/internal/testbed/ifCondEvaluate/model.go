//+build metago

// evaluate `if v, ok := mf.Value().(Foo); [ok] { ... }` case.

package ifCondEvaluate

import (
	"fmt"
	"time"

	"github.com/vvakame/til/go/metago"
)

type Foo struct {
	ID        int64
	Name      string `json:"nickname"`
	CreatedAt time.Time
}

func (obj *Foo) MarshalJSON() ([]byte, error) {
	mv := metago.ValueOf(obj)
	for _, mf := range mv.Fields() {

		if v, ok := mf.Value().(int64); ok {
			fmt.Println("a", v)
		} else {
			fmt.Println("b", v)
		}

		if v, ok := mf.Value().(int64); !ok {
			fmt.Println("a", v)
		} else {
			fmt.Println("b", v)
		}

		if v, ok := mf.Value().(int64); ok == true {
			fmt.Println("a", v)
		} else {
			fmt.Println("b", v)
		}
	}

	panic("👀")
}
