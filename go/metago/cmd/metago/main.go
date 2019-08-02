package main

import (
	"fmt"
	"os"

	"github.com/vvakame/til/go/metago"
)

func main() {
	err := metago.Process("github.com/vvakame/til/go/metago/internal/testbed")
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		os.Exit(0)
	}
}
