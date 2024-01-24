package main

import (
	"encoding/json"
	"net/http"

	"github.com/Captain-Leftovers/boot-chirpy-server/cmd/helpers"
)

func (cfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if len(params.Body) > 140 {
		helpers.RespondWithError(w, http.StatusBadRequest, "Chirp is longer than 140 characters")
		return
	}

	cleaned := helpers.CensorProfanity(params.Body)

	chirp, err := cfg.DB.CreateChirp(cleaned)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, chirp)
}
