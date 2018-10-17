package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"cloud.google.com/go/cloudtasks/apiv2beta3"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/blendle/zapdriver"
	"github.com/favclip/ucon"
	"go.mercari.io/datastore"
	"go.mercari.io/datastore/boom"
	"go.mercari.io/datastore/clouddatastore"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2beta3"
)

var dsClient datastore.Client
var logger *zap.Logger

func main() {

	var err error

	if os.Getenv("LOCAL_EXEC") == "" {
		config := zapdriver.NewProductionConfig()
		logger, err = config.Build()
		if err != nil {
			log.Fatalf("Failed to create zap logger: %v", err)
		}

	} else {
		logger = zap.NewNop()
	}
	defer logger.Sync()

	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: os.Getenv("GOOGLE_CLOUD_PROJECT"),
	})
	if err != nil {
		log.Fatalf("Failed to create stackdriver exporter: %v", err)
	}
	trace.RegisterExporter(exporter)
	defer exporter.Flush()

	dsClient, err = clouddatastore.FromContext(context.Background())
	if err != nil {
		log.Fatalf("Failed to create cloud datastore client: %v", err)
	}
	defer dsClient.Close()

	handlerMain()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)

	server := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		Handler: &ochttp.Handler{
			Handler:     ucon.DefaultMux,
			Propagation: &propagation.HTTPFormat{},
		},
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	log.Printf("running...")

	// setup graceful shutdown...
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM)
	<-sigCh

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("graceful shutdown failure: %s", err)
	}
	log.Printf("graceful shutdown successfully")
}

func handlerMain() {
	ucon.Orthodox()

	// https://cloud.google.com/appengine/docs/standard/go111/how-instances-are-managed#instance_scaling
	// Automatic scaling の時は動かないはず
	ucon.HandleFunc("*", "/_ah/start", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("on /_ah/start")
		fmt.Fprint(w, "on start!")
	})
	ucon.HandleFunc("*", "/_ah/stop", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("on /_ah/stop")
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

	logger.Debug("Hi, 1")
	logger.Info("Hi, 2")

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
		logger.Debug("loop", zap.Int("index", i), zap.Int("value", v))
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

	bm := boom.FromClient(ctx, dsClient)

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		key, err := bm.Put(&Go111SampleKind{})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%s\n", key.String())
		wg.Done()
	}()
	go func() {
		key, err := bm.Put(&Go111SampleKind{})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%s\n", key.String())
		wg.Done()
	}()
	go func() {
		key, err := bm.Put(&Go111SampleKind{})
		if err != nil {
			log.Fatal(err)
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
	logger.Info("on taskClient.CreateTask", zap.Any("task", task))

	fmt.Fprintf(w, "added!")

	return nil
}

func cloudTaskExecHandler(w http.ResponseWriter, r *http.Request) error {

	logger.Info("task done!")
	fmt.Fprintf(w, "done!")

	return nil
}
