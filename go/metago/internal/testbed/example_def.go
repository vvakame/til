package testbed

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"

	"github.com/vvakame/til/go/metago"
)

type Foo struct {
	ID        int64
	Name      string `json:"nickname"`
	CreatedAt time.Time
}

type Bar struct {
	ID        string
	NickName  string
	CreatedAt time.Time
}

// MarshalJSON return JSON format binary array.
func (foo *Foo) MarshalJSON() ([]byte, error) {
	mv := metago.ValueOf(foo)
	return marshalJSONTemplate(mv)
}

func (obj *Bar) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("{")

	mv := metago.ValueOf(obj)
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

		if v := strings.SplitN(mf.StructTagGet("json"), ",", 2)[0]; v != "" {
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

		if v := strings.SplitN(mf.StructTagGet("json"), ",", 2)[0]; v != "" {
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
