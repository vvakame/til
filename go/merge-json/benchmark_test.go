package merge_json

import (
	"bytes"
	pkg_a "github.com/vvakame/til/go/merge-json/a"
	pkg_b "github.com/vvakame/til/go/merge-json/b"
	pkg_c "github.com/vvakame/til/go/merge-json/c"
	"runtime"
	"strconv"
	"testing"
)

func Benchmark_2objects(b *testing.B) {
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
				ret, err := m.f(obj...)
				if err != nil {
					b.Fatal(err)
				}
				_ = ret
			}
		})
	}
}

func Benchmark_10objects(b *testing.B) {
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
	}

	for _, m := range marshallers {
		m := m
		b.Run(m.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				ret, err := m.f(obj...)
				if err != nil {
					b.Fatal(err)
				}
				_ = ret
			}
		})
	}
}

func Benchmark_10objectsWithLongText(b *testing.B) {
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

	var longString func(digit int) string
	{
		var buf bytes.Buffer
		longString = func(digit int) string {
			buf.Reset()
			for i := 0; i < digit; i++ {
				buf.WriteString(strconv.Itoa(i % 10))
			}
			return buf.String()
		}
	}

	strLen := 1000000
	obj := []interface{}{
		&A{
			Name:    longString(strLen),
			Age:     37,
			private: true,
		},
		&B{
			B1: longString(strLen),
			B2: 2,
		},
		&C{
			C1: longString(strLen),
			C2: longString(strLen),
		},
		&D{D: longString(strLen)},
		&E{E: longString(strLen)},
		&F{F: longString(strLen)},
		&G{G: longString(strLen)},
		&H{H: longString(strLen)},
		&I{I: longString(strLen)},
		&J{J: longString(strLen)},
	}

	runtime.GC()

	for _, m := range marshallers {
		m := m
		b.Run(m.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				ret, err := m.f(obj...)
				if err != nil {
					b.Fatal(err)
				}
				_ = ret
			}
		})
	}
}
