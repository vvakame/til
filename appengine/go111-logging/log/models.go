package log

// LogEntryDeprecated provides Stackdriver LogEntry format.
// spec: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
type LogEntryDeprecated struct {
	LogName          string               `json:"logName"`
	Resource         interface{}          `json:"resource,omitempty"`
	Timestamp        interface{}          `json:"timestamp,omitempty"`
	ReceiveTimestamp interface{}          `json:"receiveTimestamp,omitempty"`
	Severity         string               `json:"severity" validate:"enum=DEFAULT|DEBUG|INFO|NOTICE|WARNING|ERROR|CRITICAL|ALERT|EMERGENCY"`
	InsertID         string               `json:"insertId,omitempty"`
	HttpRequest      *LogEntryHttpRequest `json:"httpRequest,omitempty"`
	Labels           map[string]string    `json:"labels,omitempty"`
	Operation        *LogEntryOperation   `json:"operation,omitempty"`
	Trace            string               `json:"trace,omitempty"`
	SpanID           string               `json:"spanId,omitempty"`
	TraceSampled     *bool                `json:"traceSampled,omitempty"`
	SourceLocation   interface{}          `json:"sourceLocation,omitempty"`
	TextPayload      string               `json:"textPayload,omitempty"`
	JSONPayload      interface{}          `json:"jsonPayload,omitempty"`
}

type LogEntryOperation struct {
	ID       string `json:"id,omitempty"`
	Producer string `json:"producer,omitempty"`
	First    *bool  `json:"first,omitempty"`
	Last     *bool  `json:"last,omitempty"`
}

// LogEntryHttpRequest provides HttpRequest log.
// spec: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#httprequest
type LogEntryHttpRequest struct {
	RequestMethod                  string `json:"requestMethod"`
	RequestURL                     string `json:"requestUrl"`
	RequestSize                    int64  `json:"requestSize,string,omitempty"`
	Status                         int    `json:"status"`
	ResponseSize                   int64  `json:"responseSize,string,omitempty"`
	UserAgent                      string `json:"userAgent,omitempty"`
	RemoteIP                       string `json:"remoteIp,omitempty"`
	Referer                        string `json:"referer,omitempty"`
	Latency                        string `json:"latency,omitempty"`
	CacheLookup                    *bool  `json:"cacheLookup,omitempty"`
	CacheHit                       *bool  `json:"cacheHit,omitempty"`
	CacheValidatedWithOriginServer *bool  `json:"cacheValidatedWithOriginServer,omitempty"`
	CacheFillBytes                 *int64 `json:"cacheFillBytes,string,omitempty"`
	Protocol                       string `json:"protocol"`
}

// LogEntry is nanika.
// spec: https://cloud.google.com/logging/docs/agent/configuration#special-fields
type LogEntry struct {
	Severity       string               `json:"severity" validate:"enum=DEFAULT|DEBUG|INFO|NOTICE|WARNING|ERROR|CRITICAL|ALERT|EMERGENCY"`
	HttpRequest    *LogEntryHttpRequest `json:"httpRequest,omitempty"`
	Time           string               `json:"time,omitempty"`
	Trace          string               `json:"logging.googleapis.com/trace,omitempty"`
	SpanID         string               `json:"logging.googleapis.com/spanId,omitempty"`
	Operation      *LogEntryOperation   `json:"logging.googleapis.com/operation,omitempty"`
	SourceLocation interface{}          `json:"logging.googleapis.com/sourceLocation,omitempty"`
	Message        string               `json:"message,omitempty"`
}
