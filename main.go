package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/pat"
)

const (
	defaultServerPort = 3000
)

func main() {
	addr := ":" + strconv.Itoa(defaultServerPort)
	p := pat.New()
	p.Get("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(w, false)
	})

	log.Printf("listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, p))
}
