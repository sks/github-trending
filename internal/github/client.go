package github

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Client a client to talk to github server
type Client struct {
	httpClient   HTTPClient
	token        string
	githubAPIURL string
}

// NewClient create a new client
func NewClient(httpClient HTTPClient) Client {
	githubAPIURL := os.Getenv("GITHUB_API")
	if githubAPIURL == "" {
		githubAPIURL = "https://api.github.com"
	}
	return Client{
		httpClient:   httpClient,
		token:        os.Getenv("GITHUB_TOKEN"),
		githubAPIURL: githubAPIURL,
	}
}

// GetTrendingRepos get the trending repos
func (c Client) GetTrendingRepos(p ProjectFilter) (*TrendingReposResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/search/repositories?%s", c.githubAPIURL, p.GetParams().Encode()), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	repos := &TrendingReposResponse{}
	err = json.NewDecoder(resp.Body).Decode(repos)
	if err != nil {
		return nil, err
	}
	log.Printf("Got %d trending repos", repos.TotalCount)
	return repos, nil
}
