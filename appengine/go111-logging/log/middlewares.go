package log

import (
	"bytes"
	"context"
	"fmt"
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

	RequestLogf(ctx, b.R, 0, 0, start)

	b.R = b.R.WithContext(ctx)
	b.Context = ctx

	defer func() {
		var buf bytes.Buffer
		_ = w.Header().Write(&buf)
		fmt.Println(int64(buf.Len()), buf.String())
		responseSize := int64(buf.Len()) + w.responseSize
		RequestLogf(ctx, b.R, w.status, responseSize, start)
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
