package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	newHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
	})
	return newHandler
}

func (cfg *apiConfig) RequestNum(w http.ResponseWriter, r *http.Request) {
	value := cfg.fileserverHits.Load()
	fmt.Fprintf(w, "Hits: %v", value)
}

func (cfg *apiConfig) resetNum(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
}
