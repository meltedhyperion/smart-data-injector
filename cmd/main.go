package main

import (
	"log"
	"net/http"
)

type App struct {
	srv *http.Server
}

func main() {
	app := &App{}

	initConfig(app)
	initServer(app)

	log.Default().Printf("api running on %v", app.srv.Addr)
	log.Fatal(app.srv.ListenAndServe())

}
