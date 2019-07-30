package testbed

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"
)

type ReflectFoo struct {
	ID        int64
	Name      string `json:"nickname"`
	CreatedAt time.Time
}

func TestReflectFoo(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
	obj := &ReflectFoo{
		ID:        100,
		Name:      "SingleFoo",
		CreatedAt: time.Date(2019, 7, 30, 0, 0, 0, 0, loc),
	}

	var buf bytes.Buffer
	buf.WriteString("{")

	rv := reflect.Indirect(reflect.ValueOf(obj))
	for i := 0; i < rv.NumField(); i++ {
		if i != 0 {
			buf.WriteString(",")
		}

		rft := rv.Type().Field(i)
		rfv := rv.Field(i)

		propertyName := rft.Name
		if rft.Tag != "" {
			propertyName = strings.SplitN(rft.Tag.Get("json"), ",", 2)[0]
		}

		buf.WriteString(`"`)
		buf.WriteString(propertyName)
		buf.WriteString(`":`)

		b, err := json.Marshal(rfv.Interface())
		if err != nil {
			t.Fatal(err)
		}
		buf.Write(b)
	}

	buf.WriteString("}")

	if v := buf.String(); v != `{"ID":100,"nickname":"SingleFoo","CreatedAt":"2019-07-30T00:00:00+09:00"}` {
		t.Error("unexpected: %v", v)
	}
}
