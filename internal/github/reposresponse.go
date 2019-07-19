package github

import "time"

// TrendingReposResponse trending repos response
type TrendingReposResponse struct {
	TotalCount        int    `json:"total_count"`
	IncompleteResults bool   `json:"incomplete_results"`
	Items             []Repo `json:"items"`
}

// RepoOwner denotes the owner of the repo
type RepoOwner struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
}

// Repo denotes a repo
type Repo struct {
	Name            string      `json:"name"`
	FullName        string      `json:"full_name"`
	Owner           RepoOwner   `json:"owner"`
	HTMLURL         string      `json:"html_url"`
	Description     string      `json:"description"`
	URL             string      `json:"url"`
	TagsURL         string      `json:"tags_url"`
	LanguagesURL    string      `json:"languages_url"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	PushedAt        time.Time   `json:"pushed_at"`
	Size            int         `json:"size"`
	StargazersCount int         `json:"stargazers_count"`
	WatchersCount   int         `json:"watchers_count"`
	Language        *string     `json:"language"`
	ForksCount      int         `json:"forks_count"`
	MirrorURL       interface{} `json:"mirror_url"`
	OpenIssuesCount int         `json:"open_issues_count"`
	License         RepoLicense `json:"license"`
	Forks           int         `json:"forks"`
	Watchers        int         `json:"watchers"`
	DefaultBranch   string      `json:"default_branch"`
	Score           float64     `json:"score"`
}

// RepoLicense denotes the license for the repo
type RepoLicense struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	SpdxID string `json:"spdx_id"`
	URL    string `json:"url"`
	NodeID string `json:"node_id"`
}
