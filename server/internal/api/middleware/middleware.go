package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"

	"noted.local/noted/internal/api/problem"
	"noted.local/noted/internal/platform"
)

type contextKey string

const requestIDKey contextKey = "noted.request_id"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id, _ = platform.NewID()
		}
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ID(r *http.Request) string {
	if r == nil {
		return ""
	}
	value, _ := r.Context().Value(requestIDKey).(string)
	return value
}

func Recover(logger *slog.Logger, next http.Handler) http.Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				logger.Error("panic recovered", "request_id", ID(r), "method", r.Method, "path", r.URL.Path, "panic", recovered, "stack", string(debug.Stack()))
				problem.Write(w, ID(r), problem.New(http.StatusInternalServerError, problem.InternalError, "The server could not complete the request"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func BodyLimit(maxBytes int64, next http.Handler) http.Handler {
	if maxBytes <= 0 {
		maxBytes = 5 * 1024 * 1024
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
		}
		next.ServeHTTP(w, r)
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
func (w *statusWriter) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(data)
	w.bytes += n
	return n, err
}
func (w *statusWriter) Unwrap() http.ResponseWriter { return w.ResponseWriter }

func Logging(logger *slog.Logger, next http.Handler) http.Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		sw := &statusWriter{ResponseWriter: w}
		next.ServeHTTP(sw, r)
		status := sw.status
		if status == 0 {
			status = http.StatusOK
		}
		// Deliberately log route metadata only. Request bodies and query values can contain user content.
		logger.Info("http request", "request_id", ID(r), "method", r.Method, "path", r.URL.Path, "status", status, "bytes", sw.bytes, "duration_ms", time.Since(started).Milliseconds())
	})
}

func Chain(handler http.Handler, logger *slog.Logger, maxBody int64) http.Handler {
	return RequestID(Recover(logger, BodyLimit(maxBody, Logging(logger, handler))))
}
