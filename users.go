package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/adamsma/webserver/internal/auth"
	"github.com/adamsma/webserver/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type Credentials struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

func (cfg *apiConfig) handleCreateUser(resp http.ResponseWriter, req *http.Request) {

	type response struct {
		User
	}

	decoder := json.NewDecoder(req.Body)
	params := Credentials{}
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

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(
			resp,
			http.StatusInternalServerError,
			"Unable to create create user",
			fmt.Errorf("error hashing password (%s): %s", params.Password, err),
		)

		return
	}

	user, err := cfg.db.CreateUser(
		req.Context(), 
		database.CreateUserParams{Email: params.Email, HashedPassword: hash},
	)
	if err != nil {
		respondWithError(
			resp,
			http.StatusInternalServerError,
			"Unable to create new users",
			err,
		)

		return
	}

	newUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(resp, http.StatusCreated, response{User: newUser})
}

func (cfg *apiConfig) handleLogin(resp http.ResponseWriter, req *http.Request){

	type response struct {
		User
	}

	decoder := json.NewDecoder(req.Body)
	params := Credentials{}
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

	tgtUser, err := cfg.db.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		respondWithError(
			resp, 
			http.StatusNotFound,
			"Unknown user",
			fmt.Errorf(
				"unable to retrieve user information (%s): %s", params.Email, err,
			),
		)

		return
	}

	err = auth.CheckPasswordHash(params.Password, tgtUser.HashedPassword, )
	if err != nil {
		respondWithError(
			resp, 
			http.StatusUnauthorized,
			"Invalid password or user",
			fmt.Errorf(
				"failed login attempt for (%s): %s", params.Email, err,
			),
		)

		return
	}

	activeUser := User{
		ID:        tgtUser.ID,
		CreatedAt: tgtUser.CreatedAt,
		UpdatedAt: tgtUser.UpdatedAt,
		Email:     tgtUser.Email,
	}

	respondWithJSON(resp, http.StatusOK, response{User: activeUser})
}