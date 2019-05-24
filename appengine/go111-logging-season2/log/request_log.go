package log

import (
	"context"
	"github.com/rs/xid"
	"net/http"
	"time"
)

type contextLoggerKey struct{}
type contextOperationKey struct{}

type RequestLogger interface {
	PushSeverity(severity LogSeverity)
	Finish()
}

func NewRequestLogger(ctx context.Context, r *http.Request, w http.ResponseWriter) (context.Context, RequestLogger, http.ResponseWriter) {
	rww := &responseWriterWatcher{ResponseWriter: w}

	logger := &rootLogger{
		ctx:     ctx,
		r:       r,
		w:       rww,
		startAt: time.Now(),
	}
	ctx = context.WithValue(ctx, contextLoggerKey{}, logger)

	operationID := r.Header.Get("X-Appengine-Request-Log-Id")
	if operationID == "" {
		operationID = xid.New().String()
	}

	op := &LogEntryOperation{
		ID:       operationID,
		Producer: "github.com/vvakame/til/appengine/go111-logging",
	}
	ctx = context.WithValue(ctx, contextOperationKey{}, op)

	return ctx, logger, rww
}

type rootLogger struct {
	ctx     context.Context
	r       *http.Request
	w       *responseWriterWatcher
	startAt time.Time

	severity LogSeverity
}

func (rl *rootLogger) PushSeverity(severity LogSeverity) {
	if rl.severity < severity {
		rl.severity = severity
	}
}

func (rl *rootLogger) Finish() {
	//var buf bytes.Buffer
	//_ = w.Header().Write(&buf)
	//fmt.Println(int64(buf.Len()), buf.String())
	//responseSize := int64(buf.Len()) + w.responseSize
	responseSize := rl.w.responseSize
	RequestLog(rl.ctx, rl.r, rl.severity, rl.w.status, responseSize, rl.startAt)
}

// NOTE
//   https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#httprequest
//   responseSize
//   The size of the HTTP response message sent back to the client, in bytes, including the response headers and the response body.
//  ... but we can't it

type responseWriterWatcher struct {
	http.ResponseWriter
	status       int
	responseSize int64
}

func (w *responseWriterWatcher) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseWriterWatcher) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}

	n, err := w.ResponseWriter.Write(b)
	w.responseSize += int64(n)

	return n, err
}
