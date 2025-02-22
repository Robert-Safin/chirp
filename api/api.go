package api

import (
	"chirpy/internal/database"
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	Db             *database.Queries
	Platform       string
	JwtSecret      string
	Polka          string
}

func (cfg *ApiConfig) HandlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>`,
		cfg.FileserverHits.Load())))
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func HandlerWithConfig(cfg *ApiConfig, handler func(*ApiConfig, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(cfg, w, r)
	}
}

func Reset(cfg *ApiConfig, w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		w.WriteHeader(403)
		return
	}

	err := cfg.Db.DeleteAllUsers(context.Background())

	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	return
}
