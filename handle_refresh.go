package main

import (
	"fmt"
	"net/http"

	"github.com/adamsma/webserver/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(resp http.ResponseWriter, req *http.Request) {

	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(
			resp,
			http.StatusUnauthorized,
			"Invalid credentials",
			err,
		)
		return
	}

	details, err := cfg.db.GetRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(
			resp,
			http.StatusUnauthorized,
			"Invalid refresh token",
			fmt.Errorf(
				"unable to retrieve refresh token information (%s): %s", refreshToken, err,
			),
		)

		return
	}

	if details.IsExpired || details.RevokedAt.Valid {
		respondWithError(
			resp,
			http.StatusUnauthorized,
			"Invalid refresh token",
			fmt.Errorf(
				"expired or revoked refresh token attemp: %+v", details,
			),
		)

		return
	}

	newToken, err := auth.MakeJWT(details.UserID, cfg.secret)
	if err != nil {
		respondWithError(
			resp,
			http.StatusUnauthorized,
			"Unable to generate authorization token",
			err,
		)

		return
	}

	respondWithJSON(
		resp,
		http.StatusOK,
		response{Token: newToken},
	)

}

func (cfg *apiConfig) handlerRevokeRefresh(resp http.ResponseWriter, req *http.Request) {

	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(
			resp,
			http.StatusUnauthorized,
			"Invalid credentials",
			err,
		)
		return
	}

	err = cfg.db.RevokeRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondWithError(
			resp,
			http.StatusUnauthorized,
			"Unable to revoke refresh token",
			fmt.Errorf(
				"unable to revoke refresh token (%s): %s", refreshToken, err,
			),
		)

		return
	}

	resp.WriteHeader(http.StatusNoContent)

}
