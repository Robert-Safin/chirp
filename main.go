package main

import (
	"chirpy/api"
	"chirpy/lib"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {

	cfg := lib.SetUp()

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /app/healthz", api.HandlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.HandlerMetrics)
	mux.HandleFunc("POST /admin/reset", api.HandlerWithConfig(cfg, api.Reset))
	mux.HandleFunc("POST /api/chirps", api.HandlerWithConfig(cfg, api.CreateChirp))
	mux.HandleFunc("POST /api/users", api.HandlerWithConfig(cfg, api.CreateUser))

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
