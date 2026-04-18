package main

import (
	"net/http"
)

func main() {
	s := http.NewServeMux()
	s.Handle("/", http.FileServer(http.Dir(".")))
	server := http.Server{
		Addr:    ":8080",
		Handler: s,
	}
	server.ListenAndServe()
}
