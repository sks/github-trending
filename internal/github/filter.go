package github

import (
	"fmt"
	"net/url"
	"time"
)

// ProjectFilter denotes the filter criterias used for filetring the repos
type ProjectFilter struct {
	Language string `json:"lang"`
	Since    int    `json:"since"`
}

// GetParams create parameters out the project filter
func (p ProjectFilter) GetParams() url.Values {
	params := url.Values{}
	params.Add("sort", "stars")
	params.Add("order", "desc")
	if p.Language != "" {
		params.Add("q", fmt.Sprintf("language:%s", p.Language))
	}
	params.Add("q", fmt.Sprintf("created:>%s", p.sinceDate().Format("2006-01-02")))

	return params
}

func (p ProjectFilter) sinceDate() time.Time {
	if p.Since == 0 {
		p.Since = 7
	}
	return time.Now().AddDate(0, 0, -1*p.Since)
}
