package auth

import (
	"testing"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	uuid := uuid.New()
	secretKey := "secret key"

	tokenString, err := MakeJWT(uuid, secretKey)
	if err != nil {
		t.Fatalf("Test FAIL: unexpected error %v", err)
	}

	gotID, err := ValidateJWT(tokenString, secretKey)
	if err != nil {
		t.Fatalf("Test FAIL: unexpected error %v", err)
	}

	if uuid != gotID {
		t.Errorf("Test FAIL: expect: %v, got: %v", uuid, gotID)
	}
}
