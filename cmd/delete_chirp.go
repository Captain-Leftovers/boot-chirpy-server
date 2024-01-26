package main

import (
	"net/http"
	"strconv"

	"github.com/Captain-Leftovers/boot-chirpy-server/cmd/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	chirpId := chi.URLParam(r, "chirpID")

	chirpIntID, err := strconv.Atoi(chirpId)
	if err != nil {
		helpers.RespondWithError(w, http.StatusForbidden, "Invalid ID -> not found")
		return
	}

	token := r.Header.Get("Authorization")

	if len(token) < 8 {
		helpers.RespondWithError(w, http.StatusForbidden, "Invalid token")
		return
	}

	token = token[7:]

	claims := jwt.RegisteredClaims{}

	parsedToken, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.jwtSecret), nil

	})

	if err != nil {

		helpers.RespondWithError(w, http.StatusForbidden, "Invalid token")
		return
	}

	if !parsedToken.Valid {
		helpers.RespondWithError(w, http.StatusForbidden, "Invalid token")
		return
	}

	issuer := claims.Issuer

	if issuer != "chirpy-access" {
		helpers.RespondWithError(w, http.StatusForbidden, "Invalid token issuer")
		return
	}

	userId := claims.Subject

	userIntId, err := strconv.Atoi(userId)

	if err != nil {
		helpers.RespondWithError(w, http.StatusForbidden, "Invalid token issuer")
		return
	}

	err = cfg.DB.DeleteChirp(chirpIntID, userIntId)

	if err != nil {
		helpers.RespondWithError(w, http.StatusForbidden, "Invalid token issuer")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, "Chirp deleted")
}
