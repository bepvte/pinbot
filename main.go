package main

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"os"
	"html/template"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
	"database/sql"
)

var db *sql.DB
var t *template.Template

func main() {
	var err error

	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	p(err)

	p(db.Ping())

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS mydata (myname TEXT PRIMARY KEY, myblob JSON)")
	p(err)

	t = template.Must(template.ParseGlob("templ/*"))

	go discordStart()

	r := chi.NewRouter()

	r.Use(middleware.Recoverer, middleware.Logger, middleware.StripSlashes)

	r.Get("/favicon.ico", http.NotFound)
	r.Get("/", indexHandler)
	r.Get("/{ID:\\d+}", channelHandler)
	r.Get("/reload", reloadhandler)


	log.Println("Listening on port :"+os.Getenv("PORT"))
	log.Fatalln(http.ListenAndServe(":" + os.Getenv("PORT"), r))
}