package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/meltedhyperion/smart-data-injector/util"
	"github.com/rs/cors"
)

func initConfig(app *App) {
	err := godotenv.Load()
	if err != nil {
		log.Error(err)
	}

}

func initServer(app *App) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// setup cors
	r.Use(cors.New(cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
		ExposedHeaders:   []string{"Authorization"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowOriginFunc: func(origin string) bool {
			return true
		},
	}).Handler)

	initHandler(app, r)

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "3000"
	}
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	srv := http.Server{
		Addr:    addr,
		Handler: r,
	}
	app.srv = &srv

	walkFunc := func(method string, route string, handler http.Handler, middleware ...func(http.Handler) http.Handler) error {
		fmt.Printf("\t\t%s %s\n", util.PadStringTo(method, 7), route)
		return nil
	}

	fmt.Print("\t\tRegistered Routes: \n\n")
	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Printf("Error logging routes. Err: %s\n", err.Error())
	}
}
