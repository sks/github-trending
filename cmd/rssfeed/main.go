package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/sks/github-trending/internal/github"
	"github.com/sks/github-trending/internal/rssfeed"
)

func main() {
	githubClient := github.NewClient(http.DefaultClient)
	feeder := rssfeed.NewFeeder(githubClient)
	http.HandleFunc("/feed", feeder.Serve)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting the server on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
