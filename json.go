package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(resp http.ResponseWriter, respCode int, msg string, err error) {

	if err != nil {
		log.Println(err)
	}

	if respCode > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	type errorReturn struct {
		Error string `json:"error"`
	}

	respondWithJSON(resp, respCode, errorReturn{Error: msg})

}

func respondWithJSON(resp http.ResponseWriter, respCode int, payload interface{}) {
	resp.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		resp.WriteHeader(500)
		return
	}

	resp.WriteHeader(respCode)
	resp.Write(data)
}
