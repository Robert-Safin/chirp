package lib

import (
	"chirpy/api"
	"chirpy/internal/database"
	"database/sql"
	"log"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
)

func SetUp() *api.ApiConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbURL := os.Getenv("DB_URL")
	platoform := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("failed to connect to db")
	}
	dbQueries := database.New(db)

	cfg := api.ApiConfig{
		FileserverHits: atomic.Int32{},
		Db:             dbQueries,
		Platform:       platoform,
	}

	return &cfg
}
