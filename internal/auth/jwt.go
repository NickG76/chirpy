package auth

import (
	"fmt"
	"time"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Create a JWT for a user
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
	}

	// 🛠️ Use HMAC SHA256 — symmetric key signing
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 🔐 Sign the token using the secret
	return token.SignedString([]byte(tokenSecret))
}

// Validate a JWT and return the user ID
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// 🧪 Claims pointer to decode into
	claims := &jwt.RegisteredClaims{}

	// 🔍 Parse + validate signature + expiry
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure token is HMAC signed
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid or expired token: %w", err)
	}

	// 🧾 Parse the Subject into UUID
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid subject UUID: %w", err)
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header missing")
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authHeader, prefix) {
		return "", fmt.Errorf("authorization header must start with 'Bearer '")
	}

	return strings.TrimSpace(strings.TrimPrefix(authHeader, prefix)), nil
}
