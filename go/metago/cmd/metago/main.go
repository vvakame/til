package main

import "github.com/vvakame/til/go/metago"

func main() {
	err := metago.Process("github.com/vvakame/til/go/metago/internal/testbed")
	if err != nil {
		panic(err)
	}
}
