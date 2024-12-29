package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {

	auth := strings.Split(headers.Get("Authorization"), " ")
	if len(auth) == 1 {
		return "", fmt.Errorf("no authorization header found")
	}

	if auth[0] != "ApiKey" {
		return "", fmt.Errorf("invalid authorization type: %v", auth[0])
	}

	return auth[1], nil

}
