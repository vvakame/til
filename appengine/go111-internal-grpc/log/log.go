package log

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/logging"
	"google.golang.org/genproto/googleapis/api/monitoredres"
)

var logClient *logging.Client

var ctxKey = &struct{ temp string }{}

var zone string
var zoneOnce sync.Once

func Init() (func() error, error) {
	if v := os.Getenv("LOCAL_EXEC"); v != "" {
		return func() error { return nil }, nil
	}

	var err error
	logClient, err = logging.NewClient(context.Background(), os.Getenv("GOOGLE_CLOUD_PROJECT"))
	if err != nil {
		return nil, err
	}
	return logClient.Close, nil
}

func WithContext(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, ctxKey, r)
}

func Criticalf(ctx context.Context, format string, args ...interface{}) {
	emitf(logging.Critical, ctx, format, args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	emitf(logging.Debug, ctx, format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	emitf(logging.Error, ctx, format, args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	emitf(logging.Info, ctx, format, args...)
}

func Warningf(ctx context.Context, format string, args ...interface{}) {
	emitf(logging.Warning, ctx, format, args...)
}

func emitf(severity logging.Severity, ctx context.Context, format string, args ...interface{}) {
	if v := os.Getenv("LOCAL_EXEC"); v != "" {
		log.Printf(format, args...)
		return
	}

	r, ok := ctx.Value(ctxKey).(*http.Request)
	if !ok {
		log.Fatal("unexpected context. It doesn't have *http.Request")
	}

	zoneOnce.Do(func() {
		if v := os.Getenv("GAE_ZONE"); v != "" {
			zone = v
			return
		}
		v, err := metadata.Zone()
		if err != nil {
			log.Fatalf("can't get zone metadata: %s", err.Error())
		}
		zone = v
	})

	// TODO なんとかして os.Getenv("GAE_DEPLOYMENT_ID") を記録しておきたい…

	traceContext := r.Header.Get("X-Cloud-Trace-Context")
	traceID := strings.Split(traceContext, "/")[0]
	logger := logClient.Logger("request-log")
	logger.Log(logging.Entry{
		Severity: severity,
		Payload: map[string]interface{}{
			"serviceContext": map[string]interface{}{},
			"message":        fmt.Sprintf(format, args...),
		},
		Resource: &monitoredres.MonitoredResource{
			Type: "gae_app",
			Labels: map[string]string{
				"module_id":  os.Getenv("GAE_SERVICE"),
				"project_id": os.Getenv("GOOGLE_CLOUD_PROJECT"),
				"version_id": os.Getenv("GAE_VERSION"),
				"zone":       zone,
			},
		},
		Trace: fmt.Sprintf("projects/%s/traces/%s", os.Getenv("GOOGLE_CLOUD_PROJECT"), traceID),
	})
}
