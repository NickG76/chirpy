package auth

import (
	"time"
	"net/http"

	"github.com/google/uuid"
)

type AUTHFUNC interface {
	MakeJWT(uuid.UUID, string, time.Duration) (string, error)
	ValidateJWT(string, string) (uuid.UUID, error)
	GetBearerToken(http.Header) (string, error)
	HashPassword(string) (string, error)
	CheckPasswordHash(string, string) error
	MakeRefreshToken() (string, error)
}

func New(auth AUTHFUNC) *Authorization {
	return &Authorization{auth: auth}
}
type Authorization struct {
	auth AUTHFUNC
}

func (a *Authorization) MakeJWT(id uuid.UUID, secret string, expiry time.Duration) (string, error) {
    return a.auth.MakeJWT(id, secret, expiry)
}

func (a *Authorization) ValidateJWT(token, secret string) (uuid.UUID, error) {
    return a.auth.ValidateJWT(token, secret)
}

func (a *Authorization) GetBearerToken(header http.Header) (string, error) {
    return a.auth.GetBearerToken(header)
}

func (a *Authorization) HashPassword(password string) (string, error) {
    return a.auth.HashPassword(password)
}

func (a *Authorization) CheckPasswordHash(password, hash string) error {
    return a.auth.CheckPasswordHash(password, hash)
}

func (a *Authorization) MakeRefreshToken() (string, error) {
    return a.auth.MakeRefreshToken()
}
