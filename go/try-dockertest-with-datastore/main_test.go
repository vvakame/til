package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	clds "cloud.google.com/go/datastore"
	"github.com/ory/dockertest"
	"go.mercari.io/datastore/clouddatastore"
)

func TestDoSomething(t *testing.T) {
	ctx := context.Background()
	err := DoSomething(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	startAt := time.Now()
	log.Println("start testing")

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	pool.MaxWait = 10 * time.Second

	// NOTE: execute `docker pull google/cloud-sdk:244.0.0` before running test.
	//       because dockertest doesn't have indicator.

	log.Println("before docker run", time.Now().Sub(startAt))
	startAt = time.Now()
	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "google/cloud-sdk",
			Tag:        "244.0.0",
			Cmd: []string{
				"gcloud",
				"--project=" + os.Getenv("DATASTORE_PROJECT_ID"),
				"beta",
				"emulators",
				"datastore",
				"start",
				"--host-port=0.0.0.0:8081",
				"--no-store-on-disk",
				"--consistency=1.0",
			},
			ExposedPorts: []string{
				"8081",
			},
		},
	)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	log.Println("after docker run", time.Now().Sub(startAt))
	startAt = time.Now()

	retry := 0
	err = pool.Retry(func() error {
		retry++
		log.Println("retry", retry)

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()

		port := resource.GetPort("8081/tcp")
		emulatorHost := os.Getenv("DATASTORE_EMULATOR_HOST")
		if emulatorHost != "" && (strings.HasPrefix(emulatorHost, "0.0.0.0:") || strings.HasPrefix(emulatorHost, "localhost:")) {
			host := strings.SplitN(emulatorHost, ":", 2)[0]
			emulatorHost = fmt.Sprintf("%s:%s", host, port)
			err = os.Setenv("DATASTORE_EMULATOR_HOST", emulatorHost)
			if err != nil {
				return err
			}
		}

		baseDsCli, err := clds.NewClient(ctx, os.Getenv("DATASTORE_PROJECT_ID"))
		if err != nil {
			return err
		}
		dsCli, err := clouddatastore.FromClient(ctx, baseDsCli)
		if err != nil {
			return err
		}
		q := dsCli.NewQuery("__namespace__").KeysOnly()
		_, err = dsCli.GetAll(ctx, q, nil)
		if err != nil {
			return err
		}

		SetDatastoreClient(dsCli)

		return nil
	})
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	log.Println("after connecting to emulator", time.Now().Sub(startAt))
	startAt = time.Now()

	code := m.Run()

	SetDatastoreClient(nil)

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
