package merge_json

import (
	pkg_a "github.com/vvakame/til/go/merge-json/a"
	pkg_b "github.com/vvakame/til/go/merge-json/b"
	pkg_c "github.com/vvakame/til/go/merge-json/c"
	"testing"
)

func BenchmarkMarshallers(b *testing.B) {
	marshallers := []struct {
		name string
		f    func(objs ...interface{}) ([]byte, error)
	}{
		{
			"a",
			pkg_a.MarshalFlatten,
		},
		{
			"b",
			pkg_b.MarshalFlatten,
		},
		{
			"c",
			pkg_c.MarshalFlatten,
		},
	}

	obj := []interface{}{
		&A{
			Name:    "vv",
			Age:     37,
			private: true,
		},
		&B{
			B1: "b1",
			B2: 2,
		},
	}

	for _, m := range marshallers {
		m := m
		b.Run(m.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_, err := m.f(obj...)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
