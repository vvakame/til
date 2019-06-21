package govalidtemplate

import (
	"encoding/json"
	"testing"
)

func TestData_MarshalJSON(t *testing.T) {
	obj := &Data{
		ID:   "123",
		Name: "Yukari",
		Age:  5,
	}

	b, err := json.Marshal(obj)
	if err != nil {
		t.Fatal(err)
	}

	if v := string(b); v != `{"ID":"123","Name":"Yukari","Age":5}` {
		t.Errorf("unexpected: %+v", v)
	}
}
