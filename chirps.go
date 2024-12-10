package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handleValidateChirp(resp http.ResponseWriter, req *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	type errorReturn struct {
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		resp.WriteHeader(500)
		resp.Header().Set("Content-Type", "application/json")

		respBody := errorReturn{Error: "Something went wrong"}
		dat, _ := json.Marshal(respBody)
		resp.Write(dat)
		return
	}

	if len(params.Body) > 140 {
		resp.WriteHeader(400)
		resp.Header().Set("Content-Type", "application/json")

		respBody := errorReturn{Error: "Chirp is too long"}
		dat, _ := json.Marshal(respBody)
		resp.Write(dat)
		return
	}

	type returnVal struct {
		Valid bool `json:"valid"`
	}

	resp.WriteHeader(200)
	resp.Header().Set("Content-Type", "application/json")
	respBody := returnVal{Valid: true}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		resp.WriteHeader(500)
		resp.Header().Set("Content-Type", "application/json")

		respBody := errorReturn{Error: "Something went wrong"}
		dat, _ := json.Marshal(respBody)
		resp.Write(dat)
		return
	}
	resp.Write(dat)

}
