package merge_json

import (
	"github.com/vvakame/til/go/merge-json/a"
	"github.com/vvakame/til/go/merge-json/b"
	"github.com/vvakame/til/go/merge-json/c"
	"testing"
)

func Test_eachFuncs(t *testing.T) {
	marshallers := []struct {
		name string
		f    func(objs ...interface{}) ([]byte, error)
	}{
		{
			"a",
			a.MarshalFlatten,
		},
		{
			"b",
			b.MarshalFlatten,
		},
		{
			"c",
			c.MarshalFlatten,
		},
	}
	tests := []struct {
		name string
		args []interface{}
	}{
		{
			"merge 2 structs",
			[]interface{}{
				&A{
					Name:    "vv",
					Age:     37,
					private: true,
				},
				&B{
					B1: "b1",
					B2: 2,
				},
			},
		},
		{
			"merge 3 structs",
			[]interface{}{
				&A{
					Name:    "vv",
					Age:     37,
					private: true,
				},
				&B{
					B1: "b1",
					B2: 2,
				},
				&C{
					C1: "c1",
					C2: "c2",
				},
			},
		},
		{
			"merge 6 structs",
			[]interface{}{
				&A{
					Name:    "vv",
					Age:     37,
					private: true,
				},
				&B{
					B1: "b1",
					B2: 2,
				},
				&C{
					C1: "c1",
					C2: "c2",
				},
				&D{D: "d"},
				&E{E: "e"},
				&F{F: "f"},
				&G{G: "g"},
				&H{H: "h"},
				&I{I: "i"},
				&J{J: "j"},
			},
		},
	}
	for _, tt := range tests {
		for _, m := range marshallers {
			t.Run(m.name, func(t *testing.T) {
				got, err := m.f(tt.args...)
				if err != nil {
					t.Errorf("MarshalFlatten() error = %v", err)
					return
				}
				t.Log(string(got))
			})
		}
	}
}

type A struct {
	Name    string
	Age     int
	private bool
}

func (a *A) Hello() {}

type B struct {
	B1 string
	B2 int
}

func (b *B) Bye() {}

type C struct {
	C1 string
	// Name string // can't duplicate field
	C2 string `json:"c2super"`
}

func (c *C) Wow() {}

type D struct {
	D string
}

type E struct {
	E string
}

type F struct {
	F string
}

type G struct {
	G string
}

type H struct {
	H string
}

type I struct {
	I string
}

type J struct {
	J string
}
