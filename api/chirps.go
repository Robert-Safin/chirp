package api

import (
	"chirpy/internal/database"
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
		Body    string    `json:"body"`
		User_id uuid.UUID `json:"user_id"`
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
		UserID: params.User_id,
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
