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
	"syscall"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"github.com/favclip/ucon"
	"github.com/vvakame/til/appengine/go111-sample/log"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

func main() {

	close, err := log.Init()
	if err != nil {
		rlog.Fatalf("Failed to create client: %v", err)
	}
	defer close()

	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: os.Getenv("GOOGLE_CLOUD_PROJECT"),
	})
	if err != nil {
		rlog.Fatal(err)
	}
	trace.RegisterExporter(exporter)
	defer exporter.Flush()

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
