package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "secret_password"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if hash == password {
		t.Fatalf("Expected hash to be different from plaintext password")
	}
	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("Expected no error during comparison, got %v", err)
	}
	if !match {
		t.Fatalf("Expected password to match its own hash")
	}
	match, _ = CheckPasswordHash("wrong-password", hash)
	if match {
		t.Fatalf("Expected mismatch for incorrect password, but got a match")
	}
}

func TextUniqueHashes(t *testing.T) {
	password := "same-password"

	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)
	if hash1 == hash2 {
		t.Fatalf("Expected unique salts to produce different hashes for the same password")
	}
}
