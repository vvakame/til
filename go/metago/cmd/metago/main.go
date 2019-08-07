package main

import (
	"fmt"
	"os"

	"github.com/vvakame/til/go/metago"
)

func main() {
	p, err := metago.NewProcessor()
	if err != nil {
		panic(err)
	}
	result, err := p.Process(&metago.Config{
		TargetPackages: []string{
			"github.com/vvakame/til/go/metago/internal/testbed/basic",
		},
	})
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	for _, fileResult := range result.Results {
		fmt.Println(fileResult.Package.String())
		fmt.Println(fileResult.GeneratedCode)
	}
}
