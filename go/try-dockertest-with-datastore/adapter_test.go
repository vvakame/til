package main

import "go.mercari.io/datastore"

func SetDatastoreClient(cli datastore.Client) {
	dsCli = cli
}
