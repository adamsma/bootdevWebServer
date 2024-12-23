package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/adamsma/webserver/internal/auth"
	"github.com/adamsma/webserver/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Passord   string    `json:"-"`
}

type Credentials struct {
	Email     string        `json:"email"`
	Password  string        `json:"password"`
	ExpiresIn time.Duration `json:"expires_in_seconds"`
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
			"Unable to create new user",
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

func (cfg *apiConfig) handleLogin(resp http.ResponseWriter, req *http.Request) {

	type response struct {
		AccessToken  string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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
			"Invalid email or password",
			fmt.Errorf(
				"unable to retrieve user information (%s): %s", params.Email, err,
			),
		)

		return
	}

	err = auth.CheckPasswordHash(params.Password, tgtUser.HashedPassword)
	if err != nil {
		respondWithError(
			resp,
			http.StatusUnauthorized,
			"Invalid email or password",
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

	accessToken, err := auth.MakeJWT(activeUser.ID, cfg.secret)
	if err != nil {
		respondWithError(
			resp,
			http.StatusUnauthorized,
			"Unable to generate authorization token",
			err,
		)

		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(
			resp,
			http.StatusInternalServerError,
			"Unable to generate refresh token",
			err,
		)

		return
	}

	_, err = cfg.db.CreateRefreshToken(
		req.Context(),
		database.CreateRefreshTokenParams{Token: refreshToken, UserID: activeUser.ID},
	)
	if err != nil {
		respondWithError(
			resp,
			http.StatusInternalServerError,
			"Unable to generate refresh token",
			err,
		)

		return
	}

	respondWithJSON(
		resp,
		http.StatusOK,
		response{User: activeUser, AccessToken: accessToken, RefreshToken: refreshToken},
	)
}
