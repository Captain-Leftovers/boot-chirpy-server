package main

import (
	"net/http"
	"time"

	"github.com/Captain-Leftovers/boot-chirpy-server/cmd/helpers"
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

func (cfg *apiConfig) handleRefreshRefreshToken(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Token string `json:"token"`
	}

	defer r.Body.Close()

	rawToken := r.Header.Get("Authorization")

	if len(rawToken) < 8 {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid token")
		return
	}

	token := rawToken[7:]

	claims := jwt.RegisteredClaims{}

	parsedToken, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.jwtSecret), nil

	})

	if err != nil {

		helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	if !parsedToken.Valid {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	if claims.Issuer != "chirpy-refresh" {

		helpers.RespondWithError(w, http.StatusUnauthorized, "not a refresh token")
	}

	isRevoked, err := cfg.DB.CheckIsTokenRevoked(token)

	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, err.Error())
	}

	if isRevoked {
		helpers.RespondWithError(w, http.StatusUnauthorized, "token is revoked")
	}

	userId := claims.Subject
	expiration := time.Hour

	freshAccessToken, err := cfg.generateSignedJWT(expiration, userId, "chirpy-access")

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, response{Token: freshAccessToken})

}

func (cfg *apiConfig) handleRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	rawToken := r.Header.Get("Authorization")

	if len(rawToken) < 8 {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid token")
		return
	}

	token := rawToken[7:]

	claims := jwt.RegisteredClaims{}

	parsedToken, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.jwtSecret), nil

	})

	if err != nil {

		helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	if !parsedToken.Valid {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	if claims.Issuer != "chirpy-refresh" {

		helpers.RespondWithError(w, http.StatusUnauthorized, "not a refresh token")
	}

	err = cfg.DB.RevokeAccessToken(token)

	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, err.Error())
	}

	helpers.RespondWithJSON(w, http.StatusOK, "token revoked")

}


