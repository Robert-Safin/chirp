package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func CreateUser(cfg *ApiConfig, w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Email string `json:"email"`
	}

	type responseSuccess struct {
		Id         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Email      string    `json:"email"`
	}

	params := requestBody{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)

	if err != nil {
		RespondError(w, 500, fmt.Sprintf("something went wrong %v", err))
		return
	}

	user, err := cfg.Db.CreateUser(context.Background(), params.Email)

	if err != nil {
		RespondError(w, 400, fmt.Sprintf("error creating user record: %v", err))
		return
	}

	RespondOK(w, 201, responseSuccess{
		Id:         user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
	})

}
