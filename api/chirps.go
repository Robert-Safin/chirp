package api

import (
	"chirpy/internal/database"
	"chirpy/internal/database/auth"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

func CreateChirp(cfg *ApiConfig, w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Body string `json:"body"`
		// User_id uuid.UUID `json:"user_id"`
		Token string `json:"token"`
	}

	type responseSuccess struct {
		Id         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Body       string    `json:"body"`
		User_id    string    `json:"user_id"`
	}

	params := requestBody{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)

	if err != nil {
		RespondError(w, 500, fmt.Sprintf("something went wrong %v", err))
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondError(w, 401, "Unauthorized")
		return
	}

	user_id, err := auth.ValidateJWT(token, cfg.JwtSecret)
	if err != nil {
		RespondError(w, 401, "Unauthorized")
		return
	}

	if len(params.Body) > 140 {
		RespondError(w, 400, "Chirp is too long")
		return
	}
	banned_words := []string{"kerfuffle", "sharbert", "fornax"}

	clean_body := []string{}

	for _, word := range strings.Split(params.Body, " ") {
		if slices.Contains(banned_words, strings.ToLower(word)) {
			clean_body = append(clean_body, "****")
		} else {
			clean_body = append(clean_body, word)
		}
	}

	chirp, err := cfg.Db.CreateChirp(context.Background(), database.CreateChirpParams{
		Body:   strings.Join(clean_body, " "),
		UserID: user_id,
	})

	if err != nil {
		RespondError(w, 500, fmt.Sprintf("Error creatin chirp: %v", err))
		return
	}

	RespondOK(w, 201, responseSuccess{
		Id:         chirp.ID,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		Body:       chirp.Body,
		User_id:    uuid.UUID.String(chirp.UserID),
	})
}

func GetAllChirps(cfg *ApiConfig, w http.ResponseWriter, r *http.Request) {
	type resChirp struct {
		Id         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Body       string    `json:"body"`
		User_id    string    `json:"user_id"`
	}

	chirps, err := cfg.Db.GetAllChirps(context.Background())

	if err != nil {
		RespondError(w, 500, fmt.Sprintf("Error: %v", err))
		return
	}

	var resChirps []resChirp

	for _, c := range chirps {
		resChirps = append(resChirps, resChirp{
			Id:         c.ID,
			Created_at: c.CreatedAt,
			Updated_at: c.UpdatedAt,
			Body:       c.Body,
			User_id:    uuid.UUID.String(c.UserID),
		})
	}

	RespondOK(w, 200, resChirps)
}

func GetChirpById(cfg *ApiConfig, w http.ResponseWriter, r *http.Request) {
	type resChirp struct {
		Id         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Body       string    `json:"body"`
		User_id    string    `json:"user_id"`
	}

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		RespondError(w, 400, fmt.Sprintf("Error parsing UUID: %v", err))
		return
	}

	chirp, err := cfg.Db.GetChirpByID(context.Background(), id)
	if err != nil {
		RespondError(w, 404, fmt.Sprintf("ID not found: %v", err))
		return
	}

	RespondOK(w, 200, resChirp{
		Id:         chirp.ID,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		Body:       chirp.Body,
		User_id:    uuid.UUID.String(chirp.UserID),
	})
}
