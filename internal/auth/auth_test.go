package auth

import "testing"

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
