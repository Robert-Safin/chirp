package api

import (
	"chirpy/internal/database/auth"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func UpgradeUser(cfg *ApiConfig, w http.ResponseWriter, r *http.Request) {
	key, err := auth.GetApiKey(r.Header)
	if err != nil {
		RespondError(w, 401, "Not Authorized")
		return
	}
	if key != cfg.Polka {
		RespondError(w, 401, "Not Authorized")
		return
	}
	type body struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	params := body{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)

	if err != nil {
		RespondError(w, 500, fmt.Sprintf("something went wrong %v", err))
		return
	}

	if params.Event != "user.upgraded" {
		RespondError(w, 204, "Invalid event")
		return
	}

	_, err = cfg.Db.UpgradeUser(context.Background(), params.Data.UserID)
	if err != nil {
		RespondError(w, 404, "User not found")
		return
	}

	RespondOK(w, 204, struct{}{})
}
