package main

import (
	"context"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"errors"
	"fmt"
	"github.com/favclip/ucon"
	"github.com/vvakame/til/appengine/go111-logging/log"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {

	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: os.Getenv("GOOGLE_CLOUD_PROJECT"),
	})
	if err != nil {
		panic(err)
	}
	trace.RegisterExporter(exporter)
	defer exporter.Flush()

	handlerMain()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		Handler: &ochttp.Handler{
			Handler:     ucon.DefaultMux,
			Propagation: &propagation.HTTPFormat{},
		},
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// setup graceful shutdown...
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM)
	<-sigCh

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
}

func handlerMain() {
	ucon.Middleware(log.LoggerMiddleware)
	ucon.Orthodox()

	ucon.HandleFunc("*", "/_ah/start", func(w http.ResponseWriter, r *http.Request) {
	})
	ucon.HandleFunc("*", "/_ah/stop", func(w http.ResponseWriter, r *http.Request) {
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

	log.AppLogf(ctx, "Hello, logging! 1")

	for key, value := range r.Header {
		log.AppLogf(ctx, "%s: %+v", key, value)
	}

	log.AppLogf(ctx, "Hello, logging! 2")

	w.Write([]byte("test"))
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
		fmt.Println(v)
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
