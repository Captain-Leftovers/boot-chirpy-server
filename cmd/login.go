package main

import (
	"encoding/json"
	"net/http"

	"github.com/Captain-Leftovers/boot-chirpy-server/cmd/helpers"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	type reqBody struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)

	params := reqBody{}

	err := decoder.Decode(&params)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.LoginVerification(params.Email, params.Password)

	if err != nil {
		helpers.RespondWithError(w, http.StatusUnauthorized, err.Error())
	}

	helpers.RespondWihJSON(w, http.StatusOK, user)

}
