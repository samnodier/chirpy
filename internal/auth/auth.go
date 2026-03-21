package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	params := &argon2id.Params{
		Memory:      128 * 1024,
		Iterations:  4,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
	hash, err := argon2id.CreateHash(password, params)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, fmt.Errorf("Error comparing password: %w", err)
	}
	return match, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	reqToken := headers.Get("Authorization")
	if reqToken == "" {
		return "", fmt.Errorf("Authorization header is missing")
	}
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) != 2 {
		return "", fmt.Errorf("invalid Authorization header format")
	}

	token := strings.TrimSpace(splitToken[1])
	if token == "" {
		return "", fmt.Errorf("token is missing from authorization header")
	}
	return token, nil
}

func MakeRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func GetAPIKey(headers http.Header) (string, error) {
	reqApiKey := headers.Get("Authorization")
	if reqApiKey == "" {
		return "", fmt.Errorf("API Key authorization header is missing")
	}
	const prefix = "ApiKey "
	if !strings.HasPrefix(reqApiKey, prefix) {
		return "", fmt.Errorf("invalid authorization header format")
	}
	apiKey := strings.TrimSpace(reqApiKey[len(prefix):])
	if apiKey == "" {
		return "", fmt.Errorf("api key missing from authorization header")
	}
	return apiKey, nil
}
