package rssfeed

import "github.com/sks/github-trending/internal/github"

// GithubClient a client that can genereate the github trending repos
type GithubClient interface {
	GetTrendingRepos(p github.ProjectFilter) (*github.TrendingReposResponse, error)
}
