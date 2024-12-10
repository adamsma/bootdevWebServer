package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handleValidateChirp(resp http.ResponseWriter, req *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	type returnVal struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(
			resp,
			http.StatusInternalServerError,
			"Couldn't decode parameters",
			fmt.Errorf("Error decoding parametesr: %s", err),
		)
	}

	if len(params.Body) > 140 {
		respondWithError(resp, http.StatusBadRequest, "Chirp is too long", nil)
	}

	respondWithJSON(resp, http.StatusOK, returnVal{Valid: true})

}
