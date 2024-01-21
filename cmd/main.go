package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/Captain-Leftovers/boot-chirpy-server/internal/database"
	"github.com/Captain-Leftovers/boot-chirpy-server/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {

	dbg := flag.Bool("debug", false, "Enable debug mode")

	flag.Parse()

	if *dbg {
		err := os.Remove("../internal/database/database.json")

		if err != nil {
			log.Println("Failed when deleting database.json", err)
		}

	}

	db, err := database.InitDB("../internal/database/database.json")

	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	filepathRoot := "."
	port := "8001"

	router := chi.NewRouter()

	//main router routes
	router.Handle("/app", apiCfg.middlewareHitsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	router.Handle("/app/*", apiCfg.middlewareHitsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	apiRouter := chi.NewRouter()
	router.Mount("/api", apiRouter)

	//api routes
	apiRouter.Get("/healthz", handleHealthz)
	apiRouter.HandleFunc("/reset", apiCfg.resetHitsCount)
	apiRouter.Post("/chirps", apiCfg.handleCreateChirp)
	apiRouter.Get("/chirps", apiCfg.handleGetChirps)
	apiRouter.Get("/chirps/{id}", apiCfg.handleGetChirpById)

	apiRouter.Post("/users", apiCfg.handleCreateUser)

	adminRouter := chi.NewRouter()
	router.Mount("/admin", adminRouter)

	//admin routes
	adminRouter.Get("/metrics", apiCfg.numRequests)
	corsMux := middleware.MiddlewareCors(router)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: http://localhost:%s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}
