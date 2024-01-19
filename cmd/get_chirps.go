package main

import (
	"net/http"
	"sort"

	"github.com/Captain-Leftovers/boot-chirpy-server/cmd/helpers"
	"github.com/Captain-Leftovers/boot-chirpy-server/internal/database"
)

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.DB.GetChirps()

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Couldn't get chirps")
		return
	}

	chirps = mustSortSlice(chirps)

	helpers.RespondWihJSON(w, http.StatusOK, chirps)
}

func mustSortSlice(slice []database.Chirp) []database.Chirp {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].Id < slice[j].Id
	})
	return slice
}
