package main

import (
	"net/http"
	"os"

	"github.com/sks/github-trending/internal/github"
	"github.com/sks/github-trending/internal/rssfeed"
)

func main() {
	githubClient := github.NewClient(http.DefaultClient)
	feeder := rssfeed.NewFeeder(githubClient)
	d, err := feeder.CreateRSSFeeds(os.Getenv("REQUEST_LINK"))
	if err != nil {
		os.Stdout.WriteString(err.Error())
		os.Exit(1)
	}
	os.Stdout.WriteString(d)
}
