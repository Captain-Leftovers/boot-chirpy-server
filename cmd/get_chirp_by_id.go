package main

import (
	"net/http"
	"strconv"

	"github.com/Captain-Leftovers/boot-chirpy-server/cmd/helpers"
	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handleGetChirpById(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	idString := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, "Invalid ID -> not found")
		return
	}

	chirpById, err := cfg.DB.GetChirpById(id)

	if err != nil {
		helpers.RespondWithError(w, http.StatusNotFound, "Couldn't get chirps")
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, chirpById)
}
