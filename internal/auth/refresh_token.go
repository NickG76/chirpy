package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func MakeRefreshToken() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	n, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("Failed to generate randomn bytes: %w", err)
	}
	if n != len(bytes) {
		return "", fmt.Errorf("incomplete random bytes read")
	}
	return hex.EncodeToString(bytes), nil
}
