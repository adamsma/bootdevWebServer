package auth

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {

	pass1 := "correctPassword9!"
	pass2 := "anotherPassword934?"
	hash1, _ := HashPassword(pass1)
	hash2, _ := HashPassword(pass2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct Password",
			password: pass1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect Password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Hashes Don't Collide",
			password: pass1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty Password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid Hash",
			password: pass1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := CheckPasswordHash(test.password, test.hash)
			if (err != nil) != test.wantErr {
				t.Errorf("Check Hash error == %v, expected %v", err, test.wantErr)
			}
		})
	}

}

func TestTokenCreateAndValidate(t *testing.T) {

	id1 := uuid.New()
	id2 := uuid.New()
	secret1 := "realsecret"
	secret2 := "anotherSecret"
	futureExpire, _ := time.ParseDuration("1h")
	negExpire, _ := time.ParseDuration("-1h")

	var tokenId1Sec1, expiredToken, tokenId2Sec1, tokenId1Sec2 string

	createTests := []struct {
		name      string
		id        uuid.UUID
		secret    string
		expiry    time.Duration
		saveToken *string
		wantErr   bool
	}{
		{
			name:      "Create Valid Token",
			id:        id1,
			secret:    secret1,
			expiry:    futureExpire,
			saveToken: &tokenId1Sec1,
			wantErr:   false,
		},
		{
			name:      "Create Valid Token 2",
			id:        id2,
			secret:    secret1,
			expiry:    futureExpire,
			saveToken: &tokenId2Sec1,
			wantErr:   false,
		},
		{
			name:      "Create Valid Token 3",
			id:        id1,
			secret:    secret2,
			expiry:    futureExpire,
			saveToken: &tokenId1Sec2,
			wantErr:   false,
		},
		{
			name:      "Create Expired Token",
			id:        id1,
			secret:    secret1,
			expiry:    negExpire,
			saveToken: &expiredToken,
			wantErr:   false,
		},
	}

	for _, test := range createTests {
		t.Run(test.name, func(t *testing.T) {
			token, err := MakeJWT(test.id, test.secret, test.expiry)
			if (err != nil) != test.wantErr {
				t.Errorf("Make JWT error == %v, expected %v", err, test.wantErr)
			}

			*test.saveToken = token

		})
	}

	// test different secrets and different ids create differnt JWTs
	t.Run("Different IDs, Different JWT", func(t *testing.T) {

		if tokenId1Sec1 == tokenId2Sec1 {
			t.Errorf("different userIDs didn't produce different JWTs\ntoken1: %s\n\ntoken2:%s", tokenId1Sec1, tokenId2Sec1)
		}
	})

	t.Run("Different Secrets, Different JWT", func(t *testing.T) {
		if tokenId1Sec1 == tokenId1Sec2 {
			t.Errorf("different secrets didn't produce different JWTs\ntoken1: %s\n\ntoken2:%s", tokenId1Sec1, tokenId1Sec2)
		}
	})

	validateTests := []struct {
		name        string
		tokenString string
		secret      string
		id          uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Pass Valid Token",
			tokenString: tokenId1Sec1,
			secret:      secret1,
			id:          id1,
			wantErr:     false,
		},
		{
			name:        "Pass Valid Token 2",
			tokenString: tokenId2Sec1,
			secret:      secret1,
			id:          id2,
			wantErr:     false,
		},
		{
			name:        "Pass Valid Token 3",
			tokenString: tokenId1Sec2,
			secret:      secret2,
			id:          id1,
			wantErr:     false,
		},
		{
			name:        "Expired Token",
			tokenString: expiredToken,
			secret:      secret1,
			id:          id1,
			wantErr:     true,
		},
		{
			name:        "Incorrect Secret",
			tokenString: tokenId1Sec1,
			secret:      secret2,
			id:          id1,
			wantErr:     true,
		},
	}

	for _, test := range validateTests {
		t.Run(test.name, func(t *testing.T) {
			// t.Logf("Test %d Atempting to validate token %v", i, test.tokenString)
			id, err := ValidateJWT(test.tokenString, test.secret)
			if (err != nil) != test.wantErr {
				t.Errorf("Validate JWT error == %v, expected %v", err, test.wantErr)
			}

			if id != test.id && !test.wantErr {
				t.Errorf("expected user ID: %v,  actual user ID: %v", test.id, id)
			}
		})
	}

	// test same ID's are returned with different secrets
	t.Run("Different JWTs, Same ID", func(t *testing.T) {
		idJWT1, _ := ValidateJWT(tokenId1Sec1, secret1)
		idJWT2, _ := ValidateJWT(tokenId1Sec2, secret2)

		if idJWT1 != idJWT2 {
			t.Errorf("user ID claims should match - id1: %v, id2: %v", idJWT1, idJWT2)
		}
	})

}

func TestGetBearerToken(t *testing.T) {

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Error("error in creating new http request")
	}

	t.Run("Missing Authorization", func(t *testing.T) {
		_, err := GetBearerToken(req.Header)
		if !strings.Contains(fmt.Sprint(err), "no authorization header found") {
			t.Error("missing header  should result in error")
		}
	})

	req.Header.Set("Authorization", "Bearer token123")
	t.Run("Valid Bearer Token", func(t *testing.T) {
		token, err := GetBearerToken(req.Header)
		if err != nil {
			t.Errorf("Error extracting bearer token: %v", err)
		}
		if token != "token123" {
			t.Errorf("expected: 'token123', actual: '%s'", token)
		}
	})

	jwtExample := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiI4Mjk2MjQ3Mi04Yzk3LTQyYzktYjc3Yy1mYjcwMGY2YTE4MWYiLCJleHAiOjE3MzQ5MDQ0NzksImlhdCI6MTczNDkwMDg3OX0.9VVeJR1mB4bKp1TRxinfwy64sXoApAW7H6j5CE-TPZ8"
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtExample))
	t.Run("Valid Bearer JWT", func(t *testing.T) {
		token, err := GetBearerToken(req.Header)
		if err != nil {
			t.Errorf("Error extracting bearer token: %v", err)
		}
		if token != jwtExample {
			t.Errorf("expected: '%s', actual: '%s'", jwtExample, token)
		}
	})

	req.Header.Set("Authorization", "Basic user:password")
	t.Run("Invalid Authorization Type", func(t *testing.T) {
		_, err := GetBearerToken(req.Header)
		if !strings.Contains(fmt.Sprint(err), "invalid authorization type") {
			t.Error("Incorrect authorization type should result in error")
		}
	})

}
