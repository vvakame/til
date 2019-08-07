//+build metago

package benchmark

import (
	"strconv"
	"strings"
	"sync"
	"time"

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

var bufferPool *sync.Pool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, 100)
	},
}
var propertyNameCache map[string]string

func (obj *FooMetago) MarshalJSON() ([]byte, error) {
	buf := bufferPool.Get().([]byte)
	if propertyNameCache == nil {
		propertyNameCache = make(map[string]string)
	}

	buf = append(buf, "{"...)

	mv := metago.ValueOf(obj)
	var i int
	for _, mf := range mv.Fields() {
		if mf.Value().(time.Time).IsZero() {
			continue
		}

		if i != 0 {
			buf = append(buf, ","...)
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

		buf = append(buf, quotedPropertyName...)
		buf = append(buf, ":"...)

		switch v := mf.Value().(type) {
		case int64:
			buf = strconv.AppendInt(buf, v, 10)
		case int:
			buf = strconv.AppendInt(buf, int64(v), 10)
		case string:
			buf = strconv.AppendQuote(buf, v)
		case time.Time:
			b, err := v.MarshalJSON()
			if err != nil {
				return nil, err
			}
			buf = append(buf, b...)
		}

		i++
	}

	buf = append(buf, "}"...)
	ret := make([]byte, len(buf))
	copy(ret, buf)
	bufferPool.Put(buf[:0])

	return ret, nil
}
