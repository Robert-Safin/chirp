package api

import (
	"chirpy/internal/database"
	"chirpy/internal/database/auth"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func CreateUser(cfg *ApiConfig, w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		RespondError(w, 500, fmt.Sprintf("something went wrong %v", err))
		return
	}

	user, err := cfg.Db.CreateUser(context.Background(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hash,
	})

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

func Login(cfg *ApiConfig, w http.ResponseWriter, r *http.Request) {
	type body struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type responseSuccess struct {
		Id           uuid.UUID `json:"id"`
		Created_at   time.Time `json:"created_at"`
		Updated_at   time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}

	params := body{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)

	if err != nil {
		RespondError(w, 500, fmt.Sprintf("something went wrong %v", err))
		return
	}

	user, err := cfg.Db.FindUserByEmail(context.Background(), params.Email)
	if err != nil {
		RespondError(w, 401, "Incorrect email or password")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		RespondError(w, 401, "Incorrect email or password")
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.JwtSecret, time.Duration(1)*time.Hour)

	if err != nil {
		RespondError(w, 500, fmt.Sprintf("something went wrong %v", err))
		return
	}

	refresh_token_value, err := auth.MakeRefreshToken()

	if err != nil {
		RespondError(w, 500, fmt.Sprintf("something went wrong %v", err))
		return
	}

	refresh_token, err := cfg.Db.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{
		Token:     refresh_token_value,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Duration(60*24) * time.Hour),
	})
	if err != nil {
		RespondError(w, 500, fmt.Sprintf("something went wrong %v", err))
		return
	}

	RespondOK(w, 200, responseSuccess{
		Id:           user.ID,
		Created_at:   user.CreatedAt,
		Updated_at:   user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refresh_token.Token,
	})
}

func Refresh(cfg *ApiConfig, w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	h, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondError(w, 401, "Not authorized")
		return
	}

	token, err := cfg.Db.GetRefreshTokenByToken(context.Background(), h)
	if err != nil {
		RespondError(w, 401, "Not authorized")
		return
	}

	if token.ExpiresAt.Before(time.Now()) {
		RespondError(w, 401, "Not authorized")
		return
	}

	if token.RevokedAt.Valid {
		RespondError(w, 401, "Token has been revoked")
		return
	}

	jwt, err := auth.MakeJWT(token.UserID, cfg.JwtSecret, time.Duration(1)*time.Hour)

	if err != nil {
		RespondError(w, 500, "Something went wrong")
		return
	}

	RespondOK(w, 200, response{
		Token: jwt,
	})
}

func Revoke(cfg *ApiConfig, w http.ResponseWriter, r *http.Request) {
	h, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondError(w, 401, "Not Authorized")
		return
	}

	_, err = cfg.Db.RevokeRefreshTokenByToken(context.Background(), h)
	if err != nil {
		RespondError(w, 401, "Not Authorized")
		return
	}

	RespondOK(w, 204, struct{}{})

}
