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

type rootLogger struct{}

func (rl *rootLogger) write(ctx context.Context, r *http.Request, w *responseWriterWatcher, startAt time.Time) {
	//var buf bytes.Buffer
	//_ = w.Header().Write(&buf)
	//fmt.Println(int64(buf.Len()), buf.String())
	//responseSize := int64(buf.Len()) + w.responseSize
	responseSize := w.responseSize
	RequestLogf(ctx, r, w.status, responseSize, startAt)
}

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

	falseV := false
	httpRequestEntry := &LogEntryHttpRequest{
		RequestMethod:                  r.Method,
		RequestURL:                     u.RequestURI(),
		RequestSize:                    r.ContentLength,
		Status:                         status,
		ResponseSize:                   responseSize,
		UserAgent:                      r.UserAgent(),
		RemoteIP:                       remoteIP,
		Referer:                        r.Referer(),
		Latency:                        fmt.Sprintf("%fs", duration.Seconds()),
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
