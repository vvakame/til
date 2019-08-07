package benchmark

import (
	"encoding/json"
	"time"
)

type FooVanillaHandImpl struct {
	ID        int64
	Kind      string
	Name      string `json:"nickname"`
	Age       int
	CreatedAt time.Time
}

func (obj *FooVanillaHandImpl) MarshalJSON() ([]byte, error) {
	type X struct {
		ID        int64
		Kind      string
		Name      string `json:"nickname"`
		Age       int
		CreatedAt time.Time
	}
	x := X{
		ID:        obj.ID,
		Kind:      obj.Kind,
		Name:      obj.Name,
		Age:       obj.Age,
		CreatedAt: obj.CreatedAt,
	}
	return json.Marshal(x)
}
