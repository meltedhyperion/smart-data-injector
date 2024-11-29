package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func initHandler(app *App, r *chi.Mux) {
	r.Get("/health", func(rw http.ResponseWriter, r *http.Request) {
		sendResponse(rw, 200, nil, "Server is running")
	})
}
