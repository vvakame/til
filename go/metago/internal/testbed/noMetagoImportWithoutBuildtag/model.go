// no metago import.

package noMetagoImportWithBuildtag

import (
	"time"
)

type Foo struct {
	ID        int64
	Name      string `json:"nickname"`
	CreatedAt time.Time
}

func (obj *Foo) MarshalJSON() ([]byte, error) {
	panic("foo")
}
