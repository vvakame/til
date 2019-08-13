//+build metago

// use inline template with separate package.

package separatePackage

import (
	"time"

	"github.com/vvakame/til/go/metago"
	"github.com/vvakame/til/go/metago/internal/testbed/separatePackage/separatePackageSub"
)

type Foo struct {
	ID        int64
	Name      string `json:"nickname"`
	CreatedAt time.Time
}

func (foo *Foo) MarshalJSON() ([]byte, error) {
	mv := metago.ValueOf(foo)
	return separatePackageSub.MarshalJSONTemplate(mv)
}
