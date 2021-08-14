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
			"merge 2 struct",
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
