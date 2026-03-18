package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	secret := "my-ultra-secret"
	userID := uuid.New()

	t.Run("Valid Token", func(t *testing.T) {
		token, err := MakeJWT(userID, secret, time.Hour)
		if err != nil {
			t.Fatalf("failed to make jwt: %v", err)
		}
		returnedID, err := ValidateJWT(token, secret)
		if err != nil {
			t.Fatalf("failed to validate valid jwt: %v", err)
		}
		if returnedID != userID {
			t.Errorf("expected %v, got %v", userID, returnedID)
		}
	})

	t.Run("Expired Token", func(t *testing.T) {
		token, err := MakeJWT(userID, secret, -time.Hour)
		if err != nil {

			t.Fatalf("failed to make jwt: %v", err)
		}
		_, err = ValidateJWT(token, secret)
		if err == nil {
			t.Fatal("expected error for expired token, but got none")
		}
	})

	t.Run("Wrong Secret", func(t *testing.T) {
		token, _ := MakeJWT(userID, secret, time.Hour)
		_, err := ValidateJWT(token, "wrong-secret-123")
		if err == nil {
			t.Fatal("expected error when validating with wrong secret, but got none")
		}
	})

	t.Run("Mangled Token", func(t *testing.T) {
		token, _ := MakeJWT(userID, secret, time.Hour)
		mangledToken := token + "gargabe"
		_, err := ValidateJWT(mangledToken, secret)
		if err == nil {
			t.Fatal("expected error for mangled token string")
		}
	})
}
