package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Akash0811/chirpy/internal/backend"
	"github.com/Akash0811/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	s := http.NewServeMux()

	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found")
	}

	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Failed to connect to database due to %v\n", err)
	}
	dbQueries := database.New(db)

	cfg := backend.ApiConfig{
		FileserverHits: atomic.Int32{},
		Queries:        dbQueries,
		Platform:       platform,
		JWTSecret:      jwtSecret,
	}

	// s.Handle("/app/", http.FileServer(http.Dir(".")))
	s.Handle("/app/", http.StripPrefix("/app", cfg.MiddlewareMetricsInc(http.FileServer(http.Dir(".")))))
	s.HandleFunc("GET /api/healthz", backend.HealthzHandler)
	s.HandleFunc("GET /admin/metrics", cfg.Metrics)
	s.HandleFunc("POST /admin/reset", cfg.Reset)
	s.Handle("POST /api/validate_chirp", cfg.MiddlewareMetricsInc(http.HandlerFunc(backend.ValidateChirp)))
	s.Handle("POST /api/users", cfg.MiddlewareMetricsInc(http.HandlerFunc(cfg.AddUser)))
	s.Handle("POST /api/chirps", cfg.MiddlewareMetricsInc(http.HandlerFunc(cfg.AddChirp)))
	s.Handle("GET /api/chirps", cfg.MiddlewareMetricsInc(http.HandlerFunc(cfg.GetAllChirps)))
	s.Handle("GET /api/chirps/{chirpID}", cfg.MiddlewareMetricsInc(http.HandlerFunc(cfg.GetChirp)))
	s.Handle("POST /api/login", cfg.MiddlewareMetricsInc(http.HandlerFunc(cfg.LoginUser)))

	server := http.Server{
		Addr:    ":8080",
		Handler: s,
	}
	server.ListenAndServe()
}
