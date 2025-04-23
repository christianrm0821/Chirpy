package main

import (
	"fmt"
	"net/http"
)

// middle ware that adds 1 to the amounts of request and then serves the given http request
// Haddlerfunc is different from handlefunc
func (cfg *apiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	newHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
	return newHandler
}

// prints out the number of request made
func (cfg *apiConfig) RequestNum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	value := cfg.fileserverHits.Load()
	fmt.Fprintf(w, "<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", value)
}

// resets the number of requests made
func (cfg *apiConfig) resetNum(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
}
