package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type responseError struct {
	Error string `json:"error"`
}

func RespondError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")

	json, _ := json.Marshal(responseError{
		Error: message,
	})
	w.WriteHeader(code)
	w.Write(json)
}

func RespondOK(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)

	if err != nil {
		RespondError(w, 500, fmt.Sprintf("Error marshalling JSON: %s", err))
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}
