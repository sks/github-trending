package rssfeed

import (
	"net/http"
	"time"

	"github.com/gorilla/feeds"
	"github.com/sks/github-trending/internal/github"
)

// Feeder the feeder struct that knows how to translate the github trending repos and create rss feed out of it
type Feeder struct {
	ghClient GithubClient
}

var feed = &feeds.Feed{
	Title:       "Github Trending repos",
	Description: "Feed that lists of github trending repos",
	Author:      &feeds.Author{Name: "github-trending", Email: "noreply@gmail.com"},
	Created:     time.Now(),
	Copyright:   "See https://github.com/sks/github-trending",
}

// NewFeeder create a new RSS Feeder
func NewFeeder(ghClient GithubClient) *Feeder {
	return &Feeder{
		ghClient,
	}
}

// Serve serve the http request
func (f *Feeder) Serve(w http.ResponseWriter, r *http.Request) {
	projectFilter, err := github.NewProjectFilter(r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	fe, err := f.CreateFeeds(projectFilter, r.URL.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = feed.WriteRss(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fe.WriteRss(w)
}

// CreateRSSFeeds creates an RSS FEED
func (f *Feeder) CreateRSSFeeds(projectFilter github.ProjectFilter, link string) (string, error) {
	fe, err := f.CreateFeeds(projectFilter, link)
	if err != nil {
		return ``, err
	}
	return fe.ToRss()
}

// CreateFeeds create using the gh client
func (f *Feeder) CreateFeeds(projectFilter github.ProjectFilter, link string) (*feeds.Feed, error) {
	repos, err := f.ghClient.GetTrendingRepos(projectFilter)
	if err != nil {
		return nil, err
	}

	feed.Items, err = f.createFeeds(repos.Items)
	if err != nil {
		return nil, err
	}
	feed.Link = &feeds.Link{Href: link}
	feed.Created = time.Now()
	return feed, nil
}

func (f *Feeder) createFeeds(repos []github.Repo) (items []*feeds.Item, err error) {
	for i := range repos {
		items = append(items, &feeds.Item{
			Title:       repos[i].Name,
			Link:        &feeds.Link{Href: repos[i].HTMLURL},
			Description: repos[i].Description,
			Author:      &feeds.Author{Name: repos[i].Owner.Login},
			Created:     repos[i].CreatedAt,
		})
	}
	return
}
