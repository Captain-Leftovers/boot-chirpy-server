package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Captain-Leftovers/boot-chirpy-server/cmd/helpers"
	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

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

	helpers.RespondWihJSON(w, http.StatusCreated, publicUser)
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

	//gives error here not valid token after parsing look there

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid token")
		return
	}

	if !parsedToken.Valid {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid token")
		return
	}

	useId, err := parsedToken.Claims.GetSubject()

	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "Invalid token claim")
		return
	}

	fmt.Println(useId)

}
