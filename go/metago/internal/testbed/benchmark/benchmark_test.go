package benchmark

import (
	"encoding/json"
	"testing"
	"time"
)

func BenchmarkMetago(b *testing.B) {
	foo := &FooMetago{
		ID:        100,
		Kind:      "NFC",
		Name:      "Yukari",
		Age:       4,
		CreatedAt: time.Now(),
	}

	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(foo)
	}
}

func BenchmarkMetago_callMarshalJSON(b *testing.B) {
	foo := &FooMetago{
		ID:        100,
		Kind:      "NFC",
		Name:      "Yukari",
		Age:       4,
		CreatedAt: time.Now(),
	}

	for i := 0; i < b.N; i++ {
		_, _ = foo.MarshalJSON()
	}
}

func BenchmarkVanilla(b *testing.B) {
	foo := &FooVanilla{
		ID:        100,
		Kind:      "NFC",
		Name:      "Yukari",
		Age:       4,
		CreatedAt: time.Now(),
	}

	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(foo)
	}
}

func BenchmarkVanillaHandImpl(b *testing.B) {
	foo := &FooVanillaHandImpl{
		ID:        100,
		Kind:      "NFC",
		Name:      "Yukari",
		Age:       4,
		CreatedAt: time.Now(),
	}

	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(foo)
	}
}
