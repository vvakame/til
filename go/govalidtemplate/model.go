package govalidtemplate

import "encoding/json"

var _ json.Marshaler = (*Data)(nil)

type Data struct {
	ID   string
	Name string
	Age  int `json:"age"`
}
