//+build metago

// metago import with '.' .

package metagoDotImport

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	. "github.com/vvakame/til/go/metago"
)

type Foo struct {
	ID        int64
	Name      string `json:"nickname"`
	CreatedAt time.Time
}

func (obj *Foo) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("{")

	mv := ValueOf(obj)
	var i int
	for _, mf := range mv.Fields() {
		// continue first!
		if mf.Value().(time.Time).IsZero() {
			continue
		}

		if i != 0 {
			buf.WriteString(",")
		}

		propertyName := mf.Name()
		if v := mf.StructTagGet("json"); v != "" {
			propertyName = strings.SplitN(v, ",", 2)[0]
		}

		buf.WriteString(strconv.Quote(propertyName))
		buf.WriteString(":")

		if v, ok := mf.Value().(json.Marshaler); ok {
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
