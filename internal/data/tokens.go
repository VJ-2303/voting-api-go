package data

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	ScopeAuthentication = "authentication"
)

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

func GenerateToken(userID int64, ttl time.Duration, scope string, secret string) (*Token, error) {

	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	claims := jwt.MapClaims{
		"sub":   token.UserID,
		"scope": token.Scope,
		"exp":   jwt.NewNumericDate(token.Expiry),
		"iat":   jwt.NewNumericDate(time.Now()),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}
	token.Plaintext = signedString

	return token, nil
}
