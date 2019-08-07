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
		_, err := json.Marshal(foo)
		if err != nil {
			b.Error(err)
		}
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
		_, err := foo.MarshalJSON()
		if err != nil {
			b.Error(err)
		}
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
		_, err := json.Marshal(foo)
		if err != nil {
			b.Error(err)
		}
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
		_, err := json.Marshal(foo)
		if err != nil {
			b.Error(err)
		}
	}
}
