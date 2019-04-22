package main

import (
	"context"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"contrib.go.opencensus.io/exporter/stackdriver/propagation"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/favclip/ucon"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type LogEntry2 struct {
	Severity       string             `json:"severity" validate:"enum=DEFAULT|DEBUG|INFO|NOTICE|WARNING|ERROR|CRITICAL|ALERT|EMERGENCY"`
	Time           string             `json:"time,omitempty"`
	Trace          string             `json:"logging.googleapis.com/trace,omitempty"`
	SpanID         string             `json:"logging.googleapis.com/spanId,omitempty"`
	Operation      *LogEntryOperation `json:"logging.googleapis.com/operation,omitempty"`
	SourceLocation interface{}        `json:"logging.googleapis.com/sourceLocation,omitempty"`
	Message        string             `json:"message,omitempty"`
}

type LogEntryOperation struct {
	ID       string `json:"id,omitempty"`
	Producer string `json:"producer,omitempty"`
	First    *bool  `json:"first,omitempty"`
	Last     *bool  `json:"last,omitempty"`
}

type LogEntryTimestamp struct {
	Seconds int64 `json:"seconds,string"`
	Nanos   int64 `json:"nanos,string"`
}

type LogEntry struct {
	LogName          string               `json:"logName"`
	Resource         interface{}          `json:"resource,omitempty"`
	Timestamp        interface{}          `json:"timestamp,omitempty"`
	ReceiveTimestamp interface{}          `json:"receiveTimestamp,omitempty"`
	Severity         string               `json:"severity" validate:"enum=DEFAULT|DEBUG|INFO|NOTICE|WARNING|ERROR|CRITICAL|ALERT|EMERGENCY"`
	InsertID         string               `json:"insertId,omitempty"`
	HttpRequest      *LogEntryHttpRequest `json:"httpRequest,omitempty"`
	Labels           map[string]string    `json:"labels,omitempty"`
	Operation        interface{}          `json:"operation,omitempty"`
	Trace            string               `json:"trace,omitempty"`
	SpanID           string               `json:"spanId,omitempty"`
	TraceSampled     *bool                `json:"traceSampled,omitempty"`
	SourceLocation   interface{}          `json:"sourceLocation,omitempty"`
	TextPayload      string               `json:"textPayload,omitempty"`
	JSONPayload      interface{}          `json:"jsonPayload,omitempty"`
}

// LogEntryHttpRequest provides HttpRequest log.
// spec: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#httprequest
type LogEntryHttpRequest struct {
	RequestMethod                  string `json:"requestMethod"`
	RequestURL                     string `json:"requestUrl"`
	RequestSize                    int64  `json:"requestSize,string,omitempty"`
	Status                         int    `json:"status"`
	ResponseSize                   int64  `json:"responseSize,string,omitempty"`
	UserAgent                      string `json:"userAgent"`
	RemoteIP                       string `json:"remoteIp"`
	Referer                        string `json:"referer"`
	Latency                        string `json:"responseSize,string,omitempty"` // protobuf Duration
	CacheLookup                    *bool  `json:"cacheLookup,omitempty"`
	CacheHit                       *bool  `json:"cacheHit,omitempty"`
	CacheValidatedWithOriginServer *bool  `json:"cacheValidatedWithOriginServer,omitempty"`
	CacheFillBytes                 int64  `json:"cacheFillBytes,string,omitempty"`
	Protocol                       string `json:"protocol"`
}

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

	log2f(ctx, r, "Hello, logging! 1")

	for key, value := range r.Header {
		log2f(ctx, r, "%s: %+v", key, value)
	}

	log2f(ctx, r, "Hello, logging! 2")
}

func log2f(ctx context.Context, r *http.Request, format string, a ...interface{}) {

	traceValue := ""
	spanID := ""
	if traceHeader := r.Header.Get("X-Cloud-Trace-Context"); traceHeader != "" {
		ss := strings.SplitN(traceHeader, "/", 2)
		traceValue = fmt.Sprintf("projects/%s/traces/%s", os.Getenv("GOOGLE_CLOUD_PROJECT"), ss[0])

		if len(ss) == 2 {
			ss = strings.SplitN(ss[1], ";", 2)
			spanID = ss[0]
		}
	}

	logEntry := &LogEntry2{
		Severity: "WARNING",
		Time:     time.Now().Format(time.RFC3339Nano),
		Trace:    traceValue,
		SpanID:   spanID,
		// NOTE Operation はなくてもちゃんとグルーピングされるぽい
		Operation: &LogEntryOperation{
			ID:       r.Header.Get("X-Appengine-Request-Log-Id"),
			Producer: "appengine.googleapis.com/request_id",
		},
		Message: fmt.Sprintf(format, a...),
	}
	b, err := json.Marshal(logEntry)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func logf(ctx context.Context, r *http.Request, format string, a ...interface{}) {
	u := *r.URL
	u.Fragment = ""

	traceValue := ""
	spanID := ""
	if traceHeader := r.Header.Get("X-Cloud-Trace-Context"); traceHeader != "" {
		ss := strings.SplitN(traceHeader, "/", 2)
		traceValue = fmt.Sprintf("projects/%s/traces/%s", os.Getenv("GOOGLE_CLOUD_PROJECT"), ss[0])

		if len(ss) == 2 {
			ss = strings.SplitN(ss[1], ";", 2)
			spanID = ss[0]
		}
	}
	httpRequestEntry := &LogEntryHttpRequest{
		RequestMethod: r.Method,
		RequestURL:    u.String(),
		UserAgent:     r.UserAgent(),
		Referer:       r.Referer(),
		Protocol:      r.Proto,
	}
	fmt.Println(httpRequestEntry)

	logEntry := &LogEntry{
		LogName:  fmt.Sprintf("projects/%s/logs/test-log", os.Getenv("GOOGLE_CLOUD_PROJECT")),
		Severity: "DEBUG",
		// HttpRequest: httpRequestEntry,
		Trace:       traceValue,
		SpanID:      spanID,
		TextPayload: fmt.Sprintf(format, a...),
	}
	b, err := json.Marshal(logEntry)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
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
