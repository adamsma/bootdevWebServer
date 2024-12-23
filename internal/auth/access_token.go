package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	// TokenTypeAccess -
	TokenTypeAccess TokenType = "chirpy-access"
)

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {

	claims := &jwt.RegisteredClaims{
		Issuer:    string(TokenTypeAccess),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("unable to sign token: %w", err)
	}

	return signed, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token: %w", err)
	}

	id, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	}

	return uuid.Parse(id)
}

func GetBearerToken(headers http.Header) (string, error) {

	auth := strings.Split(headers.Get("Authorization"), " ")
	if len(auth) == 1 {
		return "", fmt.Errorf("no authorization header found")
	}

	if auth[0] != "Bearer" {
		return "", fmt.Errorf("invalid authorization type: %v", auth[0])
	}

	return auth[1], nil
}
