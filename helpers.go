package main

import (
	"io"
	"log"
	"net/http"
)

func check(err error) bool {
	if err != nil {
		log.Println(err)
		return true
	}
	return false
}

func failed(w http.ResponseWriter, err error) {
	if err != nil {
		io.WriteString(w, http.StatusText(500)+": "+err.Error())
	}
}

func p(err error) {
	if err != nil {
		panic(err)
	}
}
