package main

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"os"
	"github.com/garyburd/redigo/redis"
	"html/template"
	"github.com/go-chi/chi/middleware"
)

var db redis.Conn
var t *template.Template

func main() {
	var err error

	db, err = redis.DialURL(os.Getenv("REDIS_URL"))
	p(err)


	t = template.Must(template.ParseGlob("templ/*"))

	discordStart()

	r := chi.NewRouter()

	r.Use(middleware.Recoverer, middleware.Logger, middleware.StripSlashes)

	r.Get("/", indexHandler)
	r.Get("/{ID}", channelHandler)


	log.Println("Listening on port :"+os.Getenv("PORT"))
	log.Fatalln(http.ListenAndServe(":" + os.Getenv("PORT"), r))
}