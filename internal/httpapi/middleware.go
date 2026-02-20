package handler

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/MyNameIsWhaaat/shortener/internal/logger"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapper := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapper, r)

		logger.Info("Request",
			"method", r.Method,
			"uri", r.RequestURI,
			"remote_addr", r.RemoteAddr,
			"status_code", wrapper.statusCode,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered", "error", err, "stack", string(debug.Stack()))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}
