package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/vvakame/til/go/metago"
)

var (
	verbose = flag.Bool("v", false, "show verbose message")
)

func main() {
	flag.Parse()
	args := flag.CommandLine.Args()
	if len(args) == 0 {
		log.Fatalf("1 argument must required")
	}

	p, err := metago.NewProcessor()
	if err != nil {
		panic(err)
	}
	result, err := p.Process(&metago.Config{
		TargetPackages: args,
	})
	if err != nil {
		if result != nil {
			for _, err := range result.CompileErrors {
				_, _ = fmt.Fprintln(os.Stderr, err.Error())
			}
		} else {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
		}
		os.Exit(1)
	}

	var shouldErrorExit bool
	for _, fileResult := range result.Results {
		var fileHasError bool
		for _, nErr := range fileResult.Errors {
			switch nErr.ErrorLevel {
			case metago.ErrorLevelError, metago.ErrorLevelWarning:
				shouldErrorExit = true
				fileHasError = true
				_, _ = fmt.Fprintln(os.Stderr, nErr.Error())

			case metago.ErrorLevelNotice:
				_, _ = fmt.Fprintln(os.Stderr, nErr.Error())

			case metago.ErrorLevelDebug:
				if *verbose {
					_, _ = fmt.Fprintln(os.Stderr, nErr.Error())
				}
			}
		}
		if !fileHasError && fileResult.GeneratedCode != "" {
			if *verbose {
				var generatedFilePath string
				if cwd, err := os.Getwd(); err != nil {
					generatedFilePath = fileResult.GeneratedFilePath
				} else if p, err := filepath.Rel(cwd, fileResult.GeneratedFilePath); err != nil {
					generatedFilePath = fileResult.GeneratedFilePath
				} else {
					generatedFilePath = p
				}
				_, _ = fmt.Fprintf(os.Stderr, "generate %s\n", generatedFilePath)
			}
			err := ioutil.WriteFile(fileResult.GeneratedFilePath, []byte(fileResult.GeneratedCode), 0644)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "error on write file %s: %s", fileResult.GeneratedFilePath, err.Error())
			}
		}
	}
	if shouldErrorExit {
		os.Exit(1)
	}
}
