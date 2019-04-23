package log

import (
	"context"
	"encoding/json"
	"fmt"
	"go.opencensus.io/trace"
	"net/http"
	"os"
	"strings"
	"time"
)

func AppLogf(ctx context.Context, format string, a ...interface{}) {

	traceID := ""
	spanID := ""

	if span := trace.FromContext(ctx); span != nil {
		// X-Cloud-Trace-Context のケアはOpenCensusレベルで行っておく

		traceID = fmt.Sprintf("projects/%s/traces/%s", os.Getenv("GOOGLE_CLOUD_PROJECT"), span.SpanContext().TraceID.String())
		spanID = span.SpanContext().SpanID.String()
	}

	operation, ok := ctx.Value(contextOperationKey{}).(*LogEntryOperation)
	if !ok {
		operation = nil
	}

	logEntry := &LogEntry{
		Severity:  "INFO",
		Time:      time.Now().Format(time.RFC3339Nano),
		Trace:     traceID,
		SpanID:    spanID,
		Operation: operation,
		Message:   fmt.Sprintf(format, a...),
	}
	b, err := json.Marshal(logEntry)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func RequestLogf(ctx context.Context, r *http.Request, status int, responseSize int64, startAt time.Time) {
	u := *r.URL
	u.Fragment = ""

	traceID := ""
	spanID := ""

	if span := trace.FromContext(ctx); span != nil {
		// 一般用
		traceID = fmt.Sprintf("projects/%s/traces/%s", os.Getenv("GOOGLE_CLOUD_PROJECT"), span.SpanContext().TraceID.String())
		spanID = span.SpanContext().SpanID.String()

	} else if traceHeader := r.Header.Get("X-Cloud-Trace-Context"); traceHeader != "" {
		// AppEngine とか用
		ss := strings.SplitN(traceHeader, "/", 2)
		traceID = fmt.Sprintf("projects/%s/traces/%s", os.Getenv("GOOGLE_CLOUD_PROJECT"), ss[0])

		if len(ss) == 2 {
			ss = strings.SplitN(ss[1], ";", 2)
			spanID = ss[0]
		}
	}

	remoteIP := ""
	if v := r.Header.Get("X-AppEngine-User-IP"); v != "" {
		remoteIP = v
	} else {
		remoteIP = strings.SplitN(r.RemoteAddr, ":", 2)[0]
	}

	endAt := time.Now()
	duration := endAt.Sub(startAt)
	duration.Seconds()
	nanos := endAt.Sub(startAt).Nanoseconds()
	secs := nanos / 1e9
	nanos -= secs * 1e9

	falseV := false
	httpRequestEntry := &LogEntryHttpRequest{
		RequestMethod: r.Method,
		RequestURL:    u.RequestURI(),
		RequestSize:   r.ContentLength,
		Status:        status,
		ResponseSize:  responseSize,
		UserAgent:     r.UserAgent(),
		RemoteIP:      remoteIP,
		Referer:       r.Referer(),
		Latency: &Duration{
			Seconds: secs,
			Nanos:   int32(nanos),
		},
		CacheLookup:                    &falseV,
		CacheHit:                       &falseV,
		CacheValidatedWithOriginServer: &falseV,
		CacheFillBytes:                 nil,
		Protocol:                       r.Proto,
	}

	operation, ok := ctx.Value(contextOperationKey{}).(*LogEntryOperation)
	if !ok {
		operation = nil
	}
	first := false
	last := false
	if status == 0 {
		first = true
	} else {
		last = true
	}
	operation = &LogEntryOperation{
		ID:       operation.ID,
		Producer: operation.Producer,
		First:    &first,
		Last:     &last,
	}

	logEntry := &LogEntry{
		Severity:    "WARNING",
		Time:        endAt.Format(time.RFC3339Nano),
		HttpRequest: httpRequestEntry,
		Trace:       traceID,
		SpanID:      spanID,
		Operation:   operation,
	}
	b, err := json.Marshal(logEntry)
	if err != nil {
		panic(err)
	}
	_, _ = fmt.Fprintln(os.Stderr, string(b))
	fmt.Println("debug:", string(b))
}
