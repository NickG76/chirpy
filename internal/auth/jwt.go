package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:		"chirpy",
		Subject:	userID.String(),
		IssuedAt:   jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt:  jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		claims := jwt.RegisteredClaims{}

		token, err := jst.ParseWithClaims(tokenString, claims, func(t *jwt.token) (interface{}, error) {
			return []byte(tokenSecret), nil
		})

		if err != nil {
			return uuid.Nil, fmt.ERrorf("invalid or expired token: %w", err)
		}

		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return uuid.Nil, fmt.Errorf("invalid subject unique identifier: %w", err)
		}

		return userID, nil
	}
}
