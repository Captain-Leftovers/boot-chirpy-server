package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) generateSignedJWT(expires_in_seconds time.Duration, id string, issuer string) (string, error) {

	issuedAt := jwt.NewNumericDate(time.Now().UTC())
	expiresAt := jwt.NewNumericDate(time.Now().UTC().Add(expires_in_seconds))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    issuer,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
		Subject:   id,
	})

	// see key types for signing methods here : https://golang-jwt.github.io/jwt/usage/signing_methods/#signing-methods-and-key-types

	tokenString, err := token.SignedString([]byte(cfg.jwtSecret))

	if err != nil {
		return "", err
	}

	return tokenString, nil

}
