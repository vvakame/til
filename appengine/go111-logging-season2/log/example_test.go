package log_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/vvakame/til/appengine/go111-logging-season2/log"
	"time"
)

func Example() {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	type MyLog struct {
		*log.LogEntry
		FooBar string `json:"fooBar"`
	}
	logEntry := log.NewAppLogEntry(ctx, log.SeverityWarning)
	logEntry.Time = log.Time(time.Date(2019, 4, 23, 0, 0, 0, 0, loc))
	myLog := &MyLog{
		LogEntry: logEntry,
		FooBar:   "Hello, MyLog!",
	}

	b, err := json.Marshal(myLog)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	// Output: {"severity":"WARNING","time":"2019-04-23T00:00:00+09:00","fooBar":"Hello, MyLog!"}
}
