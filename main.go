package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func healthzHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("OK"))
}

func (cfg *apiConfig) metrics(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(http.StatusOK)
	numHits := cfg.fileserverHits.Add(0)
	resp.Write([]byte(fmt.Sprintf("Hits: %v\n", numHits)))
}

func (cfg *apiConfig) reset(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(http.StatusOK)
	cfg.fileserverHits.And(0)
	resp.Write([]byte("Hits reset to 0\n"))
}

func main() {
	s := http.NewServeMux()
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	// s.Handle("/app/", http.FileServer(http.Dir(".")))
	s.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	s.HandleFunc("/healthz", healthzHandler)
	s.HandleFunc("/metrics", cfg.metrics)
	s.HandleFunc("/reset", cfg.reset)

	server := http.Server{
		Addr:    ":8080",
		Handler: s,
	}
	server.ListenAndServe()
}
