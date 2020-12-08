package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"
)

type TestEvent struct {
	Time    time.Time // encodes as an RFC3339-format string
	Action  string
	Package string
	Test    string
	Elapsed float64 // seconds
	Output  string
}

func pipe(in io.Reader, textOut io.Writer, jsonOut io.Writer) error {
	r := bufio.NewReader(in)
	var buf bytes.Buffer
outer:
	for {
		buf.Reset()

		for {
			line, isPrefix, err := r.ReadLine()
			if errors.Is(err, io.EOF) {
				break outer
			} else if err != nil {
				return err
			}

			buf.Write(line)
			if !isPrefix {
				break
			}
		}

		_, err := jsonOut.Write(buf.Bytes())
		if err != nil {
			return err
		}
		_, err = jsonOut.Write([]byte("\n"))
		if err != nil {
			return err
		}

		obj := &TestEvent{}
		err = json.Unmarshal(buf.Bytes(), obj)
		if err != nil {
			continue
		}
		_, _ = textOut.Write([]byte(obj.Output))
	}

	return nil
}

func main() {
	f, err := os.OpenFile("testresult.nljson", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = f.Close()
	}()

	err = pipe(os.Stdin, os.Stdout, f)
	if err != nil {
		panic(err)
	}
}
