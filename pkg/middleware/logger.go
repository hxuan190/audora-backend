package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

type LogEntry struct {
	TraceID      string            `json:"trace_id"`
	SpanID       string            `json:"span_id"`
	Timestamp    time.Time         `json:"timestamp"`
	Method       string            `json:"method"`
	Path         string            `json:"path"`
	Status       int               `json:"status"`
	Latency      float64           `json:"latency_ms"`
	ClientIP     string            `json:"client_ip"`
	UserAgent    string            `json:"user_agent"`
	RequestBody  string            `json:"request_body,omitempty"`
	ResponseBody string            `json:"response_body,omitempty"`
	Headers      map[string]string `json:"headers,omitempty"`
	ServiceName  string            `json:"service_name"`
	Environment  string            `json:"environment"`
	Error        string            `json:"error,omitempty"`
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement logger - https://github.com/capy-engineer/logtrace-system
	}
}
