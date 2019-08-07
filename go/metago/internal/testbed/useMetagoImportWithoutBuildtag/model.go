// metago import without buildtag.

package useMetagoImportWithoutBuildtag

import (
	"time"

	"github.com/vvakame/til/go/metago"
)

var _ metago.Value = nil

type Foo struct {
	ID        int64
	Name      string `json:"nickname"`
	CreatedAt time.Time
}

func (obj *Foo) MarshalJSON() ([]byte, error) {
	panic("foo")
}
