package main

import (
	"net/http"
	"sync/atomic"

	"github.com/Akash0811/chirpy/internal/backend"
)

func main() {
	s := http.NewServeMux()
	cfg := backend.ApiConfig{
		FileserverHits: atomic.Int32{},
	}

	// s.Handle("/app/", http.FileServer(http.Dir(".")))
	s.Handle("/app/", http.StripPrefix("/app", cfg.MiddlewareMetricsInc(http.FileServer(http.Dir(".")))))
	s.HandleFunc("GET /api/healthz", backend.HealthzHandler)
	s.HandleFunc("GET /admin/metrics", cfg.Metrics)
	s.HandleFunc("POST /admin/reset", cfg.Reset)

	server := http.Server{
		Addr:    ":8080",
		Handler: s,
	}
	server.ListenAndServe()
}
