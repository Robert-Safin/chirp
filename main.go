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
	mux.HandleFunc("GET /api/chirps", api.HandlerWithConfig(cfg, api.GetAllChirps))
	mux.HandleFunc("GET /api/chirps/{id}", api.HandlerWithConfig(cfg, api.GetChirpById))

	mux.HandleFunc("POST /api/login", api.HandlerWithConfig(cfg, api.Login))

	mux.HandleFunc("POST /api/refresh", api.HandlerWithConfig(cfg, api.Refresh))

	mux.HandleFunc("POST /api/revoke", api.HandlerWithConfig(cfg, api.Revoke))

	mux.HandleFunc("PUT /api/users", api.HandlerWithConfig(cfg, api.UpdateUser))

	mux.HandleFunc("DELETE /api/chirps/{id}", api.HandlerWithConfig(cfg, api.DeleteChirp))

	mux.HandleFunc("POST /api/polka/webhooks", api.HandlerWithConfig(cfg, api.UpgradeUser))

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
