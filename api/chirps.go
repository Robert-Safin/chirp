package api

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func ValidateChirp(w http.ResponseWriter, r *http.Request) {
	type requestJsonContents struct {
		Body string `json:"body"`
	}

	type responseSuccess struct {
		Cleaned_body string `json:"cleaned_body"`
	}

	params := requestJsonContents{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)

	if err != nil {
		RespondError(w, 500, "Something went wrong")
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

	RespondOK(w, 200, responseSuccess{Cleaned_body: strings.Join(clean_body, " ")})

	return

}
