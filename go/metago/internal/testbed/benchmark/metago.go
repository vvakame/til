//+build metago

package benchmark

import (
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/vvakame/til/go/metago"
)

// TODO 無くても動くように
var _ metago.Value = nil

type FooMetago struct {
	ID        int64
	Kind      string
	Name      string `json:"nickname"`
	Age       int
	CreatedAt time.Time
}

var propertyNameCache map[string]string

func (obj *FooMetago) MarshalJSON() ([]byte, error) {
	var buf strings.Builder
	buf.Grow(1024)
	if propertyNameCache == nil {
		propertyNameCache = make(map[string]string)
	}

	buf.WriteString("{")

	mv := metago.ValueOf(obj)
	var i int
	for _, mf := range mv.Fields() {
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
		quotedPropertyName, ok := propertyNameCache[propertyName]
		if !ok {
			quotedPropertyName = strconv.Quote(propertyName)
			propertyNameCache[propertyName] = quotedPropertyName
		}

		buf.WriteString(quotedPropertyName)
		buf.WriteString(":")

		switch v := mf.Value().(type) {
		case int64:
			buf.WriteString(strconv.FormatInt(v, 10))
		case int:
			buf.WriteString(strconv.Itoa(v))
		case string:
			buf.WriteString(strconv.Quote(v))
		case time.Time:
			b, err := v.MarshalJSON()
			if err != nil {
				return nil, err
			}
			buf.Write(b)
		}

		i++
	}

	buf.WriteString("}")

	s := buf.String()
	return *(*[]byte)(unsafe.Pointer(&s)), nil
}
