package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Akash0811/chirpy/internal/backend"
	"github.com/Akash0811/chirpy/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	s := http.NewServeMux()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Failed to connect to database due to %v\n", err)
	}
	dbQueries := database.New(db)

	cfg := backend.ApiConfig{
		FileserverHits: atomic.Int32{},
		Queries:        dbQueries,
	}

	// s.Handle("/app/", http.FileServer(http.Dir(".")))
	s.Handle("/app/", http.StripPrefix("/app", cfg.MiddlewareMetricsInc(http.FileServer(http.Dir(".")))))
	s.HandleFunc("GET /api/healthz", backend.HealthzHandler)
	s.HandleFunc("GET /admin/metrics", cfg.Metrics)
	s.HandleFunc("POST /admin/reset", cfg.Reset)
	s.Handle("POST /api/validate_chirp", cfg.MiddlewareMetricsInc(http.HandlerFunc(backend.ValidateChirp)))

	server := http.Server{
		Addr:    ":8080",
		Handler: s,
	}
	server.ListenAndServe()
}
