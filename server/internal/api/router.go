package api

import (
	"database/sql"
	"encoding/json"
	"io/fs"
	"log/slog"
	"net/http"
	"path"
	"strings"

	"noted.local/noted/internal/api/middleware"
	"noted.local/noted/internal/api/problem"
)

type RouterOptions struct {
	DB           *sql.DB
	Version      string
	Ready        func() bool
	Assets       fs.FS
	API          http.Handler
	MaxBodyBytes int64
	Logger       *slog.Logger
}

func NewRouter(options RouterOptions) http.Handler {
	mux := http.NewServeMux()
	ready := options.Ready
	if ready == nil {
		ready = func() bool { return options.DB != nil }
	}
	version := options.Version
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok", "version": version})
	})
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		if !ready() {
			problem.Write(w, middleware.ID(r), problem.New(http.StatusServiceUnavailable, problem.ServiceUnavailable, "The service is not ready"))
			return
		}
		if options.DB != nil {
			if err := options.DB.PingContext(r.Context()); err != nil {
				problem.Write(w, middleware.ID(r), problem.New(http.StatusServiceUnavailable, problem.ServiceUnavailable, "The storage is not ready"))
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok", "version": version})
	})
	apiHandler := options.API
	if apiHandler == nil {
		apiHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			problem.Write(w, middleware.ID(r), problem.NotFoundError("API route not found"))
		})
	}
	mux.Handle("/api/v1/", apiHandler)
	mux.Handle("/api/v1", apiHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if options.Assets == nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte("<!doctype html><html><body><h1>Noted</h1></body></html>"))
			return
		}
		serveStatic(w, r, options.Assets)
	})
	return middleware.Chain(mux, options.Logger, options.MaxBodyBytes)
}

func serveStatic(w http.ResponseWriter, r *http.Request, assets fs.FS) {
	name := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
	if name == "." || name == "" {
		name = "index.html"
	}
	if file, err := fs.ReadFile(assets, name); err == nil {
		w.Header().Set("Cache-Control", cacheHeader(name))
		http.ServeContent(w, r, name, fileModTime, strings.NewReader(string(file)))
		return
	}
	for _, fallback := range []string{"200.html", "index.html", "offline.html"} {
		if file, err := fs.ReadFile(assets, fallback); err == nil {
			w.Header().Set("Cache-Control", "no-cache")
			http.ServeContent(w, r, fallback, fileModTime, strings.NewReader(string(file)))
			return
		}
	}
	http.NotFound(w, r)
}

func cacheHeader(name string) string {
	if strings.Contains(name, ".") && (strings.HasPrefix(name, "_app/immutable/") || strings.HasSuffix(name, ".woff2")) {
		return "public, max-age=31536000, immutable"
	}
	return "no-cache"
}
