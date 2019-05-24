package main

import (
	"context"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"fmt"
	"github.com/favclip/ucon"
	"github.com/vvakame/til/appengine/go111-logging-season3/log"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	exporter, err := stackdriver.NewExporter(stackdriver.Options{})
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
	ucon.Middleware(func(b *ucon.Bubble) error {
		return b.Next()
	})
	ucon.Orthodox()

	ucon.HandleFunc("*", "/_ah/start", func(w http.ResponseWriter, r *http.Request) {
	})
	ucon.HandleFunc("*", "/_ah/stop", func(w http.ResponseWriter, r *http.Request) {
	})

	ucon.HandleFunc("GET", "/count100", count100Handler)
	ucon.HandleFunc("GET", "/write1g", write1gHandler)
	ucon.HandleFunc("GET", "/", indexHandler)
}

// indexHandler responds to requests with our greeting.
func indexHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx, span := trace.StartSpan(ctx, "indexHandler")
	defer span.End()

	log.Debugf(ctx, "Hello, logging! 1")
	log.Infof(ctx, "Hello, logging! 2")
	log.Warningf(ctx, "Hello, logging! 3")

	_, _ = w.Write([]byte("indexHandler"))
}

func count100Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, span := trace.StartSpan(ctx, "count100Handler")
	defer span.End()

	for i := 1; i <= 100; i++ {
		log.Debugf(ctx, "Count %d", i)
		time.Sleep(30 * time.Millisecond)
	}

	_, _ = w.Write([]byte("count100Handler"))
}

func write1gHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, span := trace.StartSpan(ctx, "write1gHandler")
	defer span.End()

	// /var/log もtmpfsなのでは？という疑問
	// ざっくり1GB書けば死ぬでしょ(256mb limit超え)
	// message以外の部分や複数ファイルに出力している部分も多いので誤差誤差の誤差だけど

	// https://gcpug.slack.com/archives/C0D60LCAE/p1558595557080100
	// tmpfsじゃなくてoverlayfsだよって教えてもらうの巻
	// 実際F1インスタンスで900MB程度は書けた(リクエストがtimeoutして終わった

	for i := 0; i < 1024*1024; i++ {
		log.Debugf(ctx, "%d:%s", i, strings.Repeat("0123456789", 100))
	}

	_, _ = w.Write([]byte("write1gHandler"))
}
