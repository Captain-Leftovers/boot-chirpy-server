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
		Id           int    `json:"id"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	defer r.Body.Close()

	type ReqBody struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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

	accessTokenExpiration := time.Duration(1) * time.Hour
	refreshTokenExpiration := time.Duration(60) * 24 * time.Hour

	idString := fmt.Sprintf("%d", user.Id)

	jwtToken, err := cfg.generateSignedJWT(accessTokenExpiration, idString, "chirpy-access")

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refreshToken, err := cfg.generateSignedJWT(refreshTokenExpiration, idString, "chirpy-refresh")

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, UserLoginResponse{
		Id:           user.Id,
		Email:        user.Email,
		Token:        jwtToken,
		RefreshToken: refreshToken,
	})

}
