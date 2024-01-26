package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Captain-Leftovers/boot-chirpy-server/cmd/helpers"
)

func (cfg *apiConfig) handlePolkaWebhook(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	polkaApiKey := os.Getenv("POLKA_API_KEY")

	senderApiKey := r.Header.Get(("Authorization"))

	if len(senderApiKey) < 8 {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	senderApiKey = senderApiKey[7:]

	if polkaApiKey != senderApiKey {
		helpers.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			ID int `json:"user_id"`
		} `json:"data"`
	}

	isRed := false

	params := parameters{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)

	if err != nil {
		return
	}

	if params.Event != "user.upgraded" {

		helpers.RespondWithJSON(w, http.StatusOK, "not user upgraded webhook")
		return
	}

	isRed = true

	userId := params.Data.ID

	_, err = cfg.DB.UpdateUser("", "", userId, &isRed)

	if err != nil {
		fmt.Println(isRed)
		helpers.RespondWithError(w, http.StatusNotFound, "user not found")
	}

	helpers.RespondWithJSON(w, http.StatusOK, "")

}
