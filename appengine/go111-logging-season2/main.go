package main

import (
	"cloud.google.com/go/logging"
	"context"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/favclip/ucon"
	"github.com/vvakame/til/appengine/go111-logging-season2/log"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"google.golang.org/appengine"
	"google.golang.org/genproto/googleapis/api/monitoredres"
	logging2 "google.golang.org/genproto/googleapis/logging/v2"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"syscall"
	"time"
)

func main() {

	{
		files := []string{
			"/etc/google-fluentd/config.d/syslog.conf",
			"/etc/google-fluentd/config.d/forward.conf",
			"/etc/google-fluentd/google-fluentd.conf",
		}
		for _, file := range files {
			b, err := ioutil.ReadFile(file)
			if os.IsNotExist(err) {
				fmt.Println(file + " is not exists")
			} else if err != nil {
				panic(err)
			} else {
				fmt.Println(file)
				fmt.Println(string(b))
			}
		}
	}

	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID: log.GetProjectID(),
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
	ucon.HandleFunc("GET", "/count100", count100Handler)
	ucon.HandleFunc("GET", "/logging", loggingHandler)
	// ucon.HandleFunc("GET", "/undocumented", undocumentedHandler)
	ucon.HandleFunc("GET", "/panic", panicHandler)
	ucon.HandleFunc("GET", "/error", errorHandler)
	ucon.HandleFunc("GET", "/incompatible", incompatibleHandler)

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

	log.AppLogf(ctx, log.SeverityDebug, "Hello, logging! 1")

	for key, value := range r.Header {
		log.AppLogf(ctx, log.SeverityDebug, "%s: %+v", key, value)
	}

	log.AppLogf(ctx, log.SeverityDebug, "Hello, logging! 2")

	for _, kv := range os.Environ() {
		log.AppLogf(ctx, log.SeverityDebug, "%s", kv)
	}

	log.AppLogf(ctx, log.SeverityInfo, "Hello, logging! 3")

	if rand.Int()%2 == 0 {
		type MyLog struct {
			*log.LogEntry
			Name string
			Note string
		}

		logEntry := &MyLog{
			LogEntry: log.NewAppLogEntry(ctx, log.SeverityWarning),
			Name:     "yukari",
			Note:     "super kawaii cat",
		}
		b, err := json.Marshal(logEntry)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
	}

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

func count100Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, span := trace.StartSpan(ctx, "count100Handler")
	defer span.End()

	for i := 1; i <= 100; i++ {
		log.AppLogf(ctx, log.SeverityDebug, "Count %d", i)
		time.Sleep(30 * time.Millisecond)
	}

	w.Write([]byte("test"))
}

func loggingHandler(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	logCli, err := logging.NewClient(ctx, log.GetProjectID())
	if err != nil {
		return err
	}

	logger := logCli.Logger("season2")
	defer logger.Flush()

	ctx, span := trace.StartSpan(ctx, "loggingHandler")
	defer span.End()

	for i := 1; i <= 100; i++ {
		myLogEntry := log.NewAppLogEntry(ctx, log.SeverityWarning)

		logger.Log(logging.Entry{
			Severity: logging.Warning,
			Payload:  fmt.Sprintf("A Count %d", i),
			Trace:    myLogEntry.Trace,
			SpanID:   myLogEntry.SpanID,
			SourceLocation: &logging2.LogEntrySourceLocation{
				File:     myLogEntry.SourceLocation.File,
				Line:     myLogEntry.SourceLocation.Line,
				Function: myLogEntry.SourceLocation.Function,
			},
		})

		logger.Log(logging.Entry{
			Severity: logging.Warning,
			Payload:  fmt.Sprintf("B Count %d", i),
			Resource: &monitoredres.MonitoredResource{
				Type: "gae_app",
				Labels: map[string]string{
					"module_id":  appengine.ModuleName(ctx),
					"project_id": appengine.AppID(ctx),
					"version_id": appengine.VersionID(ctx),
					"zone":       appengine.Datacenter(ctx),
				},
			},
			Trace:  myLogEntry.Trace,
			SpanID: myLogEntry.SpanID,
			SourceLocation: &logging2.LogEntrySourceLocation{
				File:     myLogEntry.SourceLocation.File,
				Line:     myLogEntry.SourceLocation.Line,
				Function: myLogEntry.SourceLocation.Function,
			},
		})

		time.Sleep(30 * time.Millisecond)
	}

	w.Write([]byte("test"))
	return nil
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, span := trace.StartSpan(ctx, "panicHandler")
	defer span.End()

	defer func() {
		err := recover()
		if err != nil {
			logEntry := log.NewAppLogEntry(ctx, log.SeverityCritical)
			logEntry.Message = fmt.Sprintf("panic: %v\n\n%s", err, string(debug.Stack()))
			b, err := json.Marshal(logEntry)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(b))

			type MyLog struct {
				log.LogEntry
				Other string
			}

			myLog := &MyLog{
				LogEntry: *logEntry,
				Other:    "hello!",
			}
			b, err = json.Marshal(myLog)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(b))

			panic(err)
		}
	}()

	panic("boom!")
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, span := trace.StartSpan(ctx, "errorHandler")
	defer span.End()

	// https://cloud.google.com/error-reporting/docs/formatting-error-messages

	logEntry := log.NewAppLogEntry(ctx, log.SeverityError)
	logEntry.Message = "panic: foobar\n\n" + string(debug.Stack())

	b, err := json.Marshal(logEntry)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	w.Write([]byte("test"))
}

func incompatibleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, span := trace.StartSpan(ctx, "incompatibleHandler")
	defer span.End()

	// https://cloud.google.com/error-reporting/docs/formatting-error-messages

	type MyLog1 struct {
		log.LogEntry
		Value string
	}
	type MyLog2 struct {
		log.LogEntry
		Value []int
	}
	type MyLog3 struct {
		log.LogEntry
		Value struct {
			A string
			B string
		}
	}

	{
		myLog := &MyLog1{
			LogEntry: *log.NewAppLogEntry(ctx, log.SeverityDebug),
			Value:    "string",
		}
		b, err := json.Marshal(myLog)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
	}

	{
		myLog := &MyLog2{
			LogEntry: *log.NewAppLogEntry(ctx, log.SeverityDebug),
			Value:    []int{42, 2019},
		}
		b, err := json.Marshal(myLog)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
	}

	{
		myLog := &MyLog3{
			LogEntry: *log.NewAppLogEntry(ctx, log.SeverityDebug),
			Value: struct {
				A string
				B string
			}{
				A: "A",
				B: "B",
			},
		}
		b, err := json.Marshal(myLog)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
	}

	w.Write([]byte("test"))
}
