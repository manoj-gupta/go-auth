package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/pat"
	"github.com/joho/godotenv"
	"github.com/manoj-gupta/go-auth/oauth"
	"github.com/manoj-gupta/go-auth/oauth/github"
	"github.com/manoj-gupta/go-auth/oauth/google"
)

const (
	defaultServerPort = 3000
)

func init() {
	fmt.Println("Go Auth Demo")

	// Load .env files
	// .env.local takes predence (if present)
	godotenv.Load(".env.local")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}
}

func main() {
	oauth.RegisterProvider(
		google.New(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_SECRET"), os.Getenv("GOOGLE_REDIRECT_URL"), "email", "profile"),
		github.New(os.Getenv("GITHUB_CLIENT_ID"), os.Getenv("GITHUB_SECRET"), os.Getenv("GITHUB_REDIRECT_URL"), "user"),
	)

	p := pat.New()

	// provider auth handler
	p.Get("/auth/{provider}", func(w http.ResponseWriter, r *http.Request) {
		oauth.SignInHandler(w, r)
	})

	// index page handler
	p.Get("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(w, false)
	})

	addr := ":" + strconv.Itoa(defaultServerPort)
	log.Printf("listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, p))
}
