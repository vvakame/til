package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/vvakame/til/go/har-log"
)

func main() {
	err := realMain()
	if err != nil {
		panic(err)
	}
}

func realMain() error {
	har := &harlog.Transport{}
	hc := &http.Client{
		Transport: har,
	}

	{
		_, err := hc.Get("https://blog.vvaka.me/")
		if err != nil {
			return err
		}
	}
	{
		buf := bytes.NewBufferString(`{"test":true}`)
		req, err := http.NewRequest("POST", "https://blog.vvaka.me/", buf)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		_, err = hc.Do(req)
		if err != nil {
			return err
		}
	}
	{
		vs := url.Values{}
		vs.Add("foo", "FOO")
		vs.Add("bar", "BAR")
		buf := bytes.NewBufferString(vs.Encode())
		req, err := http.NewRequest("POST", "https://blog.vvaka.me/", buf)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		_, err = hc.Do(req)
		if err != nil {
			return err
		}
	}
	{
		var buf bytes.Buffer
		mw := multipart.NewWriter(&buf)
		{
			w, err := mw.CreateFormField("test")
			if err != nil {
				return err
			}
			_, _ = w.Write([]byte("test field"))
		}
		{
			w, err := mw.CreateFormFile("file", "hello.txt")
			if err != nil {
				return err
			}
			_, _ = w.Write([]byte("Hello, world!"))
		}
		err := mw.Close()
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", "https://blog.vvaka.me/", &buf)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", mw.FormDataContentType())
		_, err = hc.Do(req)
		if err != nil {
			return err
		}
	}

	b, err := json.MarshalIndent(har.HAR(), "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	return nil
}
