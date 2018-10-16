package main

import (
	"context"
	"fmt"
	rlog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/favclip/ucon"
	"github.com/vvakame/til/appengine/go111-sample/log"
)

func main() {

	close, err := log.Init()
	if err != nil {
		rlog.Fatalf("Failed to create client: %v", err)
	}
	defer close()

	handlerMain()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		rlog.Printf("Defaulting to port %s", port)
	}

	rlog.Printf("Listening on port %s", port)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: ucon.DefaultMux,
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

	ucon.HandleFunc("GET", "/", indexHandler)
}

// indexHandler responds to requests with our greeting.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	ctx := r.Context()

	log.Debugf(ctx, "Hi, 1")
	log.Infof(ctx, "Hi, 2")

	fmt.Fprint(w, "Hello, World!")
}
