//+build metago

// use inline template with separate files.

package separateFiles

import (
	"time"

	"github.com/vvakame/til/go/metago"
)

type Foo struct {
	ID        int64
	Name      string `json:"nickname"`
	CreatedAt time.Time
}

func (foo *Foo) MarshalJSON() ([]byte, error) {
	mv := metago.ValueOf(foo)
	return marshalJSONTemplate(mv)
}
