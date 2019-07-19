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
	Description: "discussion about tech, footie, photos",
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
	repos, err := f.ghClient.GetTrendingRepos(github.ProjectFilter{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	feed.Items, err = f.createFeeds(repos.Items)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	feed.Created = time.Now()
	feed.Link = &feeds.Link{Href: r.URL.String()}
	err = feed.WriteRss(w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
