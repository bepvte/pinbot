package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"html/template"
	"log"
	"net/http"
	"os"
)

var t *template.Template

func main() {
	t = template.Must(template.ParseGlob("templ/*"))

	go discordStart()

	r := chi.NewRouter()

	r.Use(middleware.Recoverer, middleware.Logger, middleware.StripSlashes)

	r.Get("/favicon.ico", http.NotFound)
	r.Get("/", indexHandler)
	r.Get("/{ID:\\d+}", channelHandler)
	r.Get("/reload", reloadhandler)

	log.Println("Listening on port :" + os.Getenv("PORT"))
	log.Fatalln(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
