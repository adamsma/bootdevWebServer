package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"webserver/internal/database"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handleNewChirp(resp http.ResponseWriter, req *http.Request) {

	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	type response struct {
		Chirp
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(
			resp,
			http.StatusInternalServerError,
			"Couldn't decode parameters",
			fmt.Errorf("error decoding parameters: %s", err),
		)
		return
	}

	cleanedBody, err := validateChirp(resp, params.Body)
	if err != nil {
		return
	}

	chirp, err := cfg.db.CreateChirp(
		req.Context(),
		database.CreateChirpParams{Body: cleanedBody, UserID: params.UserID},
	)
	if err != nil {
		respondWithError(
			resp,
			http.StatusInternalServerError,
			"Unable to create chirp",
			nil,
		)

		return
	}

	respondWithJSON(resp, http.StatusCreated, response{
		Chirp: Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		},
	})

}

func validateChirp(resp http.ResponseWriter, body string) (string, error) {

	if len(body) > 140 {
		respondWithError(resp, http.StatusBadRequest, "Chirp is too long", nil)
		return "", fmt.Errorf("chirp too long: %d characters", len(body))
	}

	words := strings.Split(body, " ")
	for i, word := range words {
		switch strings.ToLower(word) {
		case "kerfuffle", "sharbert", "fornax":
			words[i] = "****"
		default:
			//do nothing
		}
	}

	return strings.Join(words, " "), nil

}
