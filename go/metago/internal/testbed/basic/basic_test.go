package basic

import (
	"encoding/json"
	"testing"
)

func Test(t *testing.T) {
	obj := &Foo{
		ID:   100,
		Name: "vvakame",
	}
	b, err := json.Marshal(obj)
	if err != nil {
		t.Fatal(err)
	}

	if v := string(b); v != `{"ID":100,"nickname":"vvakame"}` {
		t.Errorf("unexpcted: %v", v)
	}
}
