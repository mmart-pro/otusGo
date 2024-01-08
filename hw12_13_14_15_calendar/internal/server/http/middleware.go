package internalhttp

import (
	"net/http"
	"time"
)

type responseWriterDecorator struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriterDecorator(w http.ResponseWriter) *responseWriterDecorator {
	return &responseWriterDecorator{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (w *responseWriterDecorator) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler, logger Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wd := NewResponseWriterDecorator(w)
		next.ServeHTTP(wd, r)
		// ip method path http response_code latency [user_agent]
		logger.Debugf("%s %s %s %s %d %v [%s]", r.RemoteAddr, r.Method, r.URL.Path, r.Proto, wd.statusCode, time.Since(start), r.UserAgent())
	})
}
