package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Captain-Leftovers/boot-chirpy-server/cmd/helpers"
)

func HandleValidate_chirp(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	type returnVal struct {
		Cleaned_Body string `json:"cleaned_body"`
	}

	type parameters struct {
		Body string `json:"body"`
	}

	params := parameters{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)

	if err != nil {
		helpers.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if len(params.Body) > 140 {
		helpers.RespondWithError(w, http.StatusBadRequest, "Chirp is longer than 140 characters")
		return
	}

	result := helpers.CensorProfanity(params.Body)

	helpers.RespondWihJSON(w, http.StatusOK, returnVal{
		Cleaned_Body: result,
	})
}
