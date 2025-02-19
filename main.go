package main

import (
	"chirpy/api"
	"fmt"
	"net/http"
	"sync/atomic"
)

func main() {
	cfg := api.ApiConfig{FileserverHits: atomic.Int32{}}

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /app/healthz", api.HandlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.HandlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.HandlerReset)
	mux.HandleFunc("POST /api/validate_chirp", api.ValidateChirp)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
