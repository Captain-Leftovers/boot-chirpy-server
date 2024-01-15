package main

import (
	"log"
	"net/http"

	"github.com/Captain-Leftovers/boot-chirpy-server/internal/handlers"
	"github.com/Captain-Leftovers/boot-chirpy-server/internal/middleware"
)

func main() {

	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	filepathRoot := "."
	port := "3000"

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareHitsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("/healthz", handlers.HealthzHandler)
	mux.HandleFunc("/metrics", apiCfg.numRequests)
	mux.HandleFunc("/reset", apiCfg.resetHitsCount)

	corsMux := middleware.MiddlewareCors(mux)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}
