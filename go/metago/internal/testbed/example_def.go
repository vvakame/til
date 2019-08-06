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

		if mf.Value().(time.Time).IsZero() {
			continue
		}

		propertyName := mf.Name()
		if v := mf.StructTagGet("json"); v != "" {
			propertyName = strings.SplitN(v, ",", 2)[0]
		}

		buf.WriteString(`"`)
		buf.WriteString(propertyName)
		buf.WriteString(`":`)

		if v, ok := mf.Value().(time.Time); ok {
			// TODO 本当は .(json.Marshaler) したい isAssignable 参照
			b, err := v.MarshalJSON()
			if err != nil {
				return nil, err
			}
			buf.Write(b)
		} else {
			b, err := json.Marshal(v)
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

		if mf.Value().(time.Time).IsZero() {
			continue
		}

		propertyName := mf.Name()
		if v := mf.StructTagGet("json"); v != "" {
			propertyName = strings.SplitN(v, ",", 2)[0]
		}

		buf.WriteString(`"`)
		buf.WriteString(propertyName)
		buf.WriteString(`":`)

		switch v := mf.Value().(type) {
		case time.Time:
			b, err := v.MarshalJSON()
			if err != nil {
				return nil, err
			}
			buf.Write(b)

		case json.Marshaler:
			b, err := v.MarshalJSON()
			if err != nil {
				return nil, err
			}
			buf.Write(b)

		default:
			b, err := json.Marshal(v)
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
