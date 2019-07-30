package testbed

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"

	"github.com/vvakame/til/go/metago"
)

type Foo1 struct {
	ID        int64
	Name      string `json:"nickname"`
	CreatedAt time.Time
}

type Foo2 struct {
	ID        int64
	NickName  string
	CreatedAt time.Time
}

type Foo1And2 metago.GenericType

var _ Foo1And2 = Foo1And2(metago.Types(
	metago.TypeHint{Receiver: new(Foo1)},
	metago.TypeHint{Receiver: new(Foo2)},
))

func (obj Foo1And2) MarshalJSON() ([]byte, error) {
	mv := metago.ValueOf(obj)
	return marshalJSONTemplate(mv)
}

func marshalJSONTemplate(mv metago.Value) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("{")

	var i int
	for _, mf := range mv.Fields() {
		if i != 0 {
			buf.WriteString(",")
		}

		if mf.MatchTypeOf(metago.TypeHint{Receiver: time.Time{}}) {
			if mf.Value().(time.Time).IsZero() {
				continue
			}
		}

		propertyName := mf.Name()

		if v := strings.SplitN(mf.StructTag().Get("json"), ",", 2)[0]; v != "" {
			propertyName = v
		}

		buf.WriteString(`"`)
		buf.WriteString(propertyName)
		buf.WriteString(`":`)

		switch v := mf.Value().(type) {
		case json.Marshaler:
			b, err := v.MarshalJSON()
			if err != nil {
				return nil, err
			}
			buf.Write(b)

		default:
			b, err := json.Marshal(mf.Value())
			if err != nil {
				return nil, err
			}
			buf.Write(b)
		}
		i++
	}

	buf.WriteString("}")

	return buf.Bytes(), nil
}
