package main

import (
	"fmt"
	rlog "log"
	"net/http"
	"os"

	"github.com/favclip/ucon"
	"github.com/vvakame/til/appengine/go111-sample/log"
)

func main() {

	realMain()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		rlog.Printf("Defaulting to port %s", port)
	}

	rlog.Printf("Listening on port %s", port)
	rlog.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), ucon.DefaultMux))
}

func realMain() {
	close, err := log.Init()
	if err != nil {
		rlog.Fatalf("Failed to create client: %v", err)
	}
	defer close()

	ucon.Middleware(func(b *ucon.Bubble) error {
		b.Context = log.WithContext(b.Context, b.R)
		b.R = b.R.WithContext(b.Context)
		return b.Next()
	})
	ucon.Orthodox()

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
