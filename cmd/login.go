package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Captain-Leftovers/boot-chirpy-server/cmd/helpers"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {

	type UserLoginResponse struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
		Token string `json:"token"`
	}

	defer r.Body.Close()

	type ReqBody struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)

	params := ReqBody{}

	err := decoder.Decode(&params)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	

	user, err := cfg.DB.LoginVerification(params.Email, params.Password)

	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// expires_in_seconds is an optional parameter. If it's specified by the client, use it as the expiration time. If it's not specified, use a default expiration time of 24 hours. If the client specified a number over 24 hours, use 24 hours as the expiration time.

	expirationToSet := time.Duration(24) * time.Hour

	if params.ExpiresInSeconds != nil {
		givenExpiration := time.Duration(*params.ExpiresInSeconds) * time.Second
		if givenExpiration < expirationToSet && givenExpiration > 0 {
			expirationToSet = givenExpiration
		}
	}

	idString := fmt.Sprintf("%d", user.Id)

	jwtToken, err := cfg.generateSignedJWT(expirationToSet, idString)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.RespondWihJSON(w, http.StatusOK, UserLoginResponse{
		Id:    user.Id,
		Email: user.Email,
		Token: jwtToken,
	})

}
