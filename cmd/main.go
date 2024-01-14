package main

import (
	"net/http"

	"github.com/Captain-Leftovers/boot-chirpy-server/internal/middleware"
)

func main() {

	mux := http.NewServeMux()

	corsMux := middleware.MiddlewareCors(mux)

	server := http.Server{
		Addr:    ":3000",
		Handler: corsMux,
	}

	server.ListenAndServe()

}
