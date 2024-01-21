package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Captain-Leftovers/boot-chirpy-server/cmd/helpers"
)

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {

	log.Println("in create user handle")

	defer r.Body.Close()

	type parameters struct {
		Email string `json:"email"`
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

	user, err := cfg.DB.CreateUser(params.Email)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	helpers.RespondWihJSON(w, http.StatusCreated, user)
}
