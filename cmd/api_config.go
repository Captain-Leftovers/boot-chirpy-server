package main

import (
	"fmt"
	"net/http"

	"github.com/Captain-Leftovers/boot-chirpy-server/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func (cfg *apiConfig) middlewareHitsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cfg.fileserverHits += 1

		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) numRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-type", "text-html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>
	`, cfg.fileserverHits)))
}

func (cfg *apiConfig) resetHitsCount(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)

}
