package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/Captain-Leftovers/boot-chirpy-server/cmd/helpers"
	"github.com/Captain-Leftovers/boot-chirpy-server/internal/database"
)

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.DB.GetChirps()

	authorID := r.URL.Query().Get("author_id")

	order := r.URL.Query().Get("sort")

	if order == "desc" {
		order = "desc"
	} else {
		order = "asc"
	}

	filteredChirps := []database.Chirp{}

	if authorID != "" {
		idInt, err := strconv.Atoi(authorID)

		if err != nil {
			filteredChirps = chirps
		} else {
			for _, chirp := range chirps {
				if chirp.AuthorId == idInt {
					filteredChirps = append(filteredChirps, chirp)
				}
			}

		}

	} else {
		filteredChirps = chirps
	}

	fmt.Println(filteredChirps)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Couldn't get chirps")
		return
	}

	sortedChirps := mustSortSlice(filteredChirps, order)

	helpers.RespondWithJSON(w, http.StatusOK, sortedChirps)
}

func mustSortSlice(slice []database.Chirp, order string) []database.Chirp {
	if order == "asc" {

		sort.Slice(slice, func(i, j int) bool {
			return slice[i].Id < slice[j].Id
		})
	} else {
		sort.Slice(slice, func(i, j int) bool {
			return slice[i].Id > slice[j].Id
		})
	}
	return slice
}
