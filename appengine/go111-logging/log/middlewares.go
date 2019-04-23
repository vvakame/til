package log

import (
	"context"
	"github.com/favclip/ucon"
	"github.com/rs/xid"
	"net/http"
	"time"
)

var _ (ucon.MiddlewareFunc) = LoggerMiddleware

type contextOperationKey struct{}

type responseWriterWatcher struct {
	http.ResponseWriter
	status       int
	responseSize int64
}

func LoggerMiddleware(b *ucon.Bubble) error {
	start := time.Now()

	ctx := b.Context

	w := &responseWriterWatcher{ResponseWriter: b.W}
	b.W = w

	operationID := b.R.Header.Get("X-Appengine-Request-Log-Id")
	if operationID == "" {
		operationID = xid.New().String()
	}

	op := &LogEntryOperation{
		ID: operationID,
		// Producer: "appengine.googleapis.com/request_id",
		Producer: "github.com/vvakame/til/appengine/go111-logging",
	}
	ctx = context.WithValue(ctx, contextOperationKey{}, op)

	b.R = b.R.WithContext(ctx)
	b.Context = ctx

	rootLogger := &rootLogger{}
	defer func() {
		rootLogger.write(b.Context, b.R, w, start)
	}()

	return b.Next()
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
