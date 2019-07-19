package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/sks/github-trending/internal/github"
	"github.com/sks/github-trending/internal/rssfeed"
)

func main() {
	githubClient := github.NewClient(http.DefaultClient)
	feeder := rssfeed.NewFeeder(githubClient)
	pf := github.ProjectFilter{
		Language: os.Getenv("LANGUAGE"),
		Since:    getNumber("SINCE", 7),
	}
	d, err := feeder.CreateRSSFeeds(pf, os.Getenv("REQUEST_LINK"))
	if err != nil {
		os.Stdout.WriteString(err.Error())
		os.Exit(1)
	}
	os.Stdout.WriteString(d)
}

func getNumber(key string, defaultValue int) int {
	val := os.Getenv(key)
	i, err := strconv.Atoi(val)
	if err == nil {
		return i
	}
	return defaultValue
}
