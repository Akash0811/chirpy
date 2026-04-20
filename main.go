package main

import (
	"net/http"
)

func healthzHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("OK"))
}

func main() {
	s := http.NewServeMux()

	// s.Handle("/app/", http.FileServer(http.Dir(".")))
	s.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	s.HandleFunc(
		"/healthz",
		healthzHandler,
	)

	server := http.Server{
		Addr:    ":8080",
		Handler: s,
	}
	server.ListenAndServe()
}
