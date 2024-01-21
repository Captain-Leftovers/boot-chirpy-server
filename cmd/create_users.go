package main

import (
	"encoding/json"
	"net/http"

	"github.com/Captain-Leftovers/boot-chirpy-server/cmd/helpers"
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
