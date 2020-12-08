package testjson

import (
	"fmt"
	"testing"
	"time"
)

func TestA(t *testing.T) {
	t.Log("Hello, test!")
	fmt.Println("Hello, stdout!")

	t.Run("1st", func(t *testing.T) {
		t.Log("in 1st")
		fmt.Println("Hi, 1st")
	})
	t.Run("2nd", func(t *testing.T) {
		t.Fail()
		fmt.Println("Hi, 2nd")
	})
	t.Run("3rd", func(t *testing.T) {
		t.Parallel()
		t.Log("in 3rd")
		time.Sleep(1 * time.Second)
		fmt.Println("Hi, 3rd")
	})

	fmt.Println("Bye, stdout!")
}
