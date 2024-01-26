package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Captain-Leftovers/boot-chirpy-server/cmd/helpers"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type parametersType struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)

	params := parametersType{}

	err := decoder.Decode(&params)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if len(params.Email) > 40 {
		helpers.RespondWithError(w, http.StatusBadRequest, "Email is longer than 40 characters")
		return
	}

	if len(params.Password) > 140 {
		helpers.RespondWithError(w, http.StatusBadRequest, "Email is longer than 140 characters")
		return
	}

	publicUser, err := cfg.DB.CreateUser(params.Email, params.Password)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Couldn't create user reason: "+err.Error())
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, publicUser)
}

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {

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

	params := parametersType{}

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&params)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), 4)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	userIdInt, err := strconv.Atoi(userId)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Couldn't convert userId to int")
		return
	}

	publicUser, err := cfg.DB.UpdateUser(params.Email, string(hashedPassword), userIdInt, nil)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Couldn't update user reason: "+err.Error())
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, publicUser)

}
