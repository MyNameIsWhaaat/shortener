package handler

import (
    "log"
    "net/http"
    "runtime/debug"
    "time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        wrapper := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}
        
        next.ServeHTTP(wrapper, r)

        log.Printf(
            "[%s] %s %s %d %s",
            r.Method,
            r.RequestURI,
            r.RemoteAddr,
            wrapper.statusCode,
            time.Since(start),
        )
    })
}

func RecoverMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("panic recovered: %v\n%s", err, debug.Stack())
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