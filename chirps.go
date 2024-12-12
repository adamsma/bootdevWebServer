package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func handleValidateChirp(resp http.ResponseWriter, req *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	type returnVal struct {
		CleanedBody string `json:"cleaned_body"`
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

	if len(params.Body) > 140 {
		respondWithError(resp, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	words := strings.Split(params.Body, " ")
	for i, word := range words {
		switch strings.ToLower(word) {
		case "kerfuffle", "sharbert", "fornax":
			words[i] = "****"
		default:
			//do nothing
		}
	}

	respondWithJSON(
		resp, http.StatusOK, returnVal{CleanedBody: strings.Join(words, " ")},
	)

}
