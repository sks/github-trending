package github

import (
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/schema"
)

// ProjectFilter denotes the filter criterias used for filetring the repos
type ProjectFilter struct {
	Language string `json:"lang"`
	Since    int    `json:"since"`
}

var decoder = schema.NewDecoder()

// NewProjectFilter create new project filter from url params
func NewProjectFilter(urlParams map[string][]string) (ProjectFilter, error) {
	p := ProjectFilter{}
	err := decoder.Decode(&p, urlParams)
	return p, err
}

// GetParams create parameters out the project filter
func (p ProjectFilter) GetParams() string {
	params := "?order=desc&s=stars"
	queries := []string{
		fmt.Sprintf("created:>%s", p.sinceDate().Format("2006-01-02")),
		fmt.Sprintf("stars:>0"),
	}
	if p.Language != "" {
		queries = append(queries, fmt.Sprintf("language:%s", p.Language))
	}
	return params + "&q=" + strings.Join(queries, "+")
}

func (p ProjectFilter) sinceDate() time.Time {
	if p.Since == 0 {
		p.Since = 7
	}
	return time.Now().AddDate(0, 0, -1*p.Since)
}
