package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"cloud.google.com/go/storage"
	"github.com/vvakame/go-harlog"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

var (
	bucket = flag.String("bucket", os.Getenv("GCS_BUCKET"), "operating GCS bucket name")
)

func init() {
	flag.Parse()
}

func main() {
	if *bucket == "" {
		panic("-bucket is required")
	}

	ctx := context.Background()

	hc, err := google.DefaultClient(ctx, storage.ScopeReadWrite)
	if err != nil {
		panic(err)
	}

	// inject HAR logger!
	har := &harlog.Transport{
		Transport: hc.Transport,
	}
	hc.Transport = har

	client, err := storage.NewClient(
		ctx,
		option.WithHTTPClient(hc),
	)
	if err != nil {
		panic(err)
	}

	bucket := client.Bucket(*bucket)

	{
		object := bucket.Object("2019-10-01-harlog/hello.txt")
		r, err := object.NewReader(ctx)
		if err != nil {
			panic(err)
		}
		b, err := ioutil.ReadAll(r)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
	}
	{
		object := bucket.Object("2019-10-01-harlog/goodnight.txt")
		w := object.NewWriter(ctx)
		_, err = w.Write([]byte("Good night, world!"))
		if err != nil {
			panic(err)
		}
		err = w.Close()
		if err != nil {
			panic(err)
		}
	}

	// dump HAR file!
	b, err := json.MarshalIndent(har.HAR(), "", "  ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("gcs.har", b, 0644)
	if err != nil {
		panic(err)
	}
}
