package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/adamsma/webserver/internal/auth"
	"github.com/adamsma/webserver/internal/database"

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
		Body string `json:"body"`
	}

	type response struct {
		Chirp
	}

	authToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(
			resp,
			http.StatusUnauthorized,
			"Invalid credentials",
			err,
		)
		return
	}

	userID, err := auth.ValidateJWT(authToken, cfg.secret)
	if err != nil {
		respondWithError(
			resp,
			http.StatusUnauthorized,
			"Invalid credentials",
			err,
		)
		return
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
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
		database.CreateChirpParams{Body: cleanedBody, UserID: userID},
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

func (cfg *apiConfig) handleGetChirps(resp http.ResponseWriter, req *http.Request) {

	chirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		respondWithError(
			resp,
			http.StatusInternalServerError,
			"Unable to retrieve chirps",
			nil,
		)

		return
	}

	var returnChirps []Chirp
	for _, chirp := range chirps {
		returnChirps = append(returnChirps, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(resp, http.StatusOK, returnChirps)

}

func (cfg *apiConfig) handleGetChirpByID(resp http.ResponseWriter, req *http.Request) {

	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondWithError(
			resp,
			http.StatusBadRequest,
			"Invalid chirpID",
			nil,
		)

		return
	}

	chirp, err := cfg.db.GetChirpByID(req.Context(), chirpID)
	if err != nil {
		respondWithError(
			resp,
			http.StatusNotFound,
			"Chirp not found",
			nil,
		)

		return
	}

	respondWithJSON(resp, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
