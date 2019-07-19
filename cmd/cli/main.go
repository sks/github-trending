package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/sks/github-trending/internal/github"
)

const format = "| %s\t| %s\t\t| %s\t\t|\n"

func main() {
	githubClient := github.NewClient(http.DefaultClient)
	resp, err := githubClient.GetTrendingRepos(github.ProjectFilter{})
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	const padding = 1
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)

	items := resp.Items[1:11]
	fmt.Fprintf(w, format, " ", "NAME", "URL")
	fmt.Fprintf(w, format, " ", "---", "---")
	for i := range items {
		fmt.Fprintf(w, format, strconv.Itoa(i), items[i].Name, items[i].URL)
	}
	w.Flush()
}
