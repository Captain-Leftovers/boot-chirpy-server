package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Captain-Leftovers/boot-chirpy-server/cmd/helpers"
	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	token := r.Header.Get("Authorization")

	if len(token) < 8 {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid token")
		return
	}

	token = token[7:]

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

	issuer := claims.Issuer

	if issuer != "chirpy-access" {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid token issuer")
		return
	}

	userId := claims.Subject

	userIdInt, err := strconv.Atoi(userId)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Couldn't convert user id to int")
		return
	}

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err = decoder.Decode(&params)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if len(params.Body) > 140 {
		helpers.RespondWithError(w, http.StatusBadRequest, "Chirp is longer than 140 characters")
		return
	}

	cleaned := helpers.CensorProfanity(params.Body)

	chirp, err := cfg.DB.CreateChirp(cleaned, userIdInt)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, chirp)

}
