package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/adamsma/webserver/internal/auth"
	"github.com/adamsma/webserver/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlePolkaWebhook(resp http.ResponseWriter, req *http.Request) {

	type parameters struct {
		Event string `json:"Event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(
			resp,
			http.StatusUnauthorized,
			"Invalid credentials",
			err,
		)
		return
	}

	if apiKey != cfg.paymentKey {
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

	if params.Event != "user.upgraded" {
		resp.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpdateChirpyRedStatus(
		req.Context(),
		database.UpdateChirpyRedStatusParams{
			ID: params.Data.UserID, IsChirpyRed: true,
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			respondWithError(
				resp,
				http.StatusNotFound,
				"Unable to find user",
				err,
			)

			return
		}

		respondWithError(resp, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	resp.WriteHeader(http.StatusNoContent)

}
