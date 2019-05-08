package main

import (
	"context"
	"go.mercari.io/datastore"
	"go.mercari.io/datastore/boom"
	"go.mercari.io/datastore/clouddatastore"
	"log"
)

var dsCli datastore.Client

func main() {
	ctx := context.Background()
	var err error
	dsCli, err = clouddatastore.FromContext(ctx)
	if err != nil {
		panic(err)
	}
	err = DoSomething(ctx)
	if err != nil {
		panic(err)
	}
}

type Something struct {
	ID   int64 `datastore:"-" boom:"id"`
	Name string
}

func DoSomething(ctx context.Context) error {
	bm := boom.FromClient(ctx, dsCli)

	obj := &Something{Name: "test"}
	_, err := bm.Put(obj)
	if err != nil {
		return err
	}
	log.Printf("%#v\n", obj)

	return nil
}
