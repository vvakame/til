package main

import (
	"context"
	"errors"
	"fmt"
	rlog "log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"cloud.google.com/go/cloudtasks/apiv2beta3"
	clouddatastore "cloud.google.com/go/datastore"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/favclip/ucon"
	"github.com/vvakame/til/appengine/go111-sample/log"
	"go.mercari.io/datastore/boom"
	cloudds "go.mercari.io/datastore/clouddatastore"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2beta3"
)

var dsClient *clouddatastore.Client

func main() {

	close, err := log.Init()
	if err != nil {
		rlog.Fatalf("Failed to create logger: %v", err)
	}
	defer close()

	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: os.Getenv("GOOGLE_CLOUD_PROJECT"),
	})
	if err != nil {
		rlog.Fatalf("Failed to create stackdriver exporter: %v", err)
	}
	trace.RegisterExporter(exporter)
	defer exporter.Flush()

	dsClient, err = clouddatastore.NewClient(context.Background(), os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		rlog.Fatalf("Failed to create cloud datastore client: %v", err)
	}
	defer dsClient.Close()

	handlerMain()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		rlog.Printf("Defaulting to port %s", port)
	}

	rlog.Printf("Listening on port %s", port)

	server := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		Handler: &ochttp.Handler{
			Handler:     ucon.DefaultMux,
			Propagation: &propagation.HTTPFormat{},
		},
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			rlog.Fatal(err)
		}
	}()

	rlog.Printf("running...")

	// setup graceful shutdown...
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM)
	<-sigCh

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err := server.Shutdown(ctx); err != nil {
		rlog.Fatalf("graceful shutdown failure: %s", err)
	}
	rlog.Printf("graceful shutdown successfully")
}

func handlerMain() {
	ucon.Middleware(func(b *ucon.Bubble) error {
		b.Context = log.WithContext(b.Context, b.R)
		b.R = b.R.WithContext(b.Context)
		return b.Next()
	})
	ucon.Orthodox()

	// https://cloud.google.com/appengine/docs/standard/go111/how-instances-are-managed#instance_scaling
	// Automatic scaling の時は動かないはず
	ucon.HandleFunc("*", "/_ah/start", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Infof(ctx, "on /_ah/start")
		fmt.Fprint(w, "on start!")
	})
	ucon.HandleFunc("*", "/_ah/stop", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Infof(ctx, "on /_ah/stop")
		fmt.Fprint(w, "on stop!")
	})

	ucon.HandleFunc("GET", "/fibonacci", fibonacciHandler)
	ucon.HandleFunc("GET", "/datastore", datastoreHandler)
	ucon.HandleFunc("GET", "/tasks", cloudTasksHandler)
	ucon.HandleFunc("GET", "/task-receive", cloudTaskExecHandler)

	ucon.HandleFunc("GET", "/", indexHandler)
}

// indexHandler responds to requests with our greeting.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	ctx := r.Context()

	ctx, span := trace.StartSpan(ctx, "indexHandler")
	defer span.End()

	log.Debugf(ctx, "Hi, 1")
	log.Infof(ctx, "Hi, 2")

	fmt.Fprint(w, "Hello, World!")
}

func fibonacciHandler(w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	ctx, span := trace.StartSpan(ctx, "fibonacciHandler")
	defer span.End()

	f := fibonacci()

	err := r.ParseForm()
	if err != nil {
		return err
	}

	v := r.Form.Get("value")
	target := 100
	if v != "" {
		target, err = strconv.Atoi(v)
		if err != nil {
			return err
		}
	}
	if 1000 <= target {
		return errors.New("bomb!!")
	}

	for i := 0; i < target; i++ {
		ctx, span = trace.StartSpan(ctx, fmt.Sprintf("fibonacciHandler#%d", i))
		span := span
		defer span.End()

		v := f()
		log.Debugf(ctx, "#%d: %d", i, v)
		fmt.Fprintf(w, "%d\n", v)
	}

	return nil
}

func fibonacci() func() int {
	f, g := 1, 0
	return func() int {
		f, g = g, f+g
		return f
	}
}

type Go111SampleKind struct {
	Kind string `datastore:"-" boom:"kind,go111-sample"`
	ID   int64  `datastore:"-" boom:"id"`
}

func datastoreHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	var bm *boom.Boom
	{
		client, err := cloudds.FromClient(ctx, dsClient)
		if err != nil {
			return err
		}
		bm = boom.FromClient(ctx, client)
	}

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		key, err := bm.Put(&Go111SampleKind{})
		if err != nil {
			rlog.Fatal(err)
		}
		fmt.Fprintf(w, "%s\n", key.String())
		wg.Done()
	}()
	go func() {
		key, err := bm.Put(&Go111SampleKind{})
		if err != nil {
			rlog.Fatal(err)
		}
		fmt.Fprintf(w, "%s\n", key.String())
		wg.Done()
	}()
	go func() {
		key, err := bm.Put(&Go111SampleKind{})
		if err != nil {
			rlog.Fatal(err)
		}
		fmt.Fprintf(w, "%s\n", key.String())
		wg.Done()
	}()
	wg.Wait()

	return nil
}

func cloudTasksHandler(w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	taskClient, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return err
	}
	task, err := taskClient.CreateTask(ctx, &taskspb.CreateTaskRequest{
		// TODO
		Parent: fmt.Sprintf("projects/%s/locations/%s/queues/go111-sample-queue", os.Getenv("GOOGLE_CLOUD_PROJECT"), "asia-northeast1"),
		Task: &taskspb.Task{
			PayloadType: &taskspb.Task_AppEngineHttpRequest{
				AppEngineHttpRequest: &taskspb.AppEngineHttpRequest{
					HttpMethod: taskspb.HttpMethod_GET,
					AppEngineRouting: &taskspb.AppEngineRouting{
						Service: os.Getenv("GAE_SERVICE"),
					},
					RelativeUri: "/task-receive",
				},
			},
		},
	})
	if err != nil {
		return err
	}
	log.Infof(ctx, "on taskClient.CreateTask: %+v", task)

	fmt.Fprintf(w, "added!")

	return nil
}

func cloudTaskExecHandler(w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	log.Infof(ctx, "task done!")
	fmt.Fprintf(w, "done!")

	return nil
}
