package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func MakeRefreshToken() (string, error) {

	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {

		return "", fmt.Errorf("error generating refresh token: %v", err)
	}

	return hex.EncodeToString(b), nil

}
