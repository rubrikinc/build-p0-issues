package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func deleteLabel(buildLabels map[string]string, name, color string) bool {

	if _, ok := buildLabels[name]; ok {

		if buildLabels[name] == color {
			return false
		}
	}

	return true

}

func createLabel(repoLabels map[string]string, name, color string) bool {

	if _, ok := repoLabels[name]; ok {

		if repoLabels[name] == color {
			return false
		}
	}

	fmt.Println(name)

	return true

}

func main() {

	ctx := context.Background()

	githubToken := os.Getenv("GITHUB_AUTH")
	// Authentication configuration
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	buildLabels := map[string]string{
		"area-cdm":                "84B6EB",
		"area-partner":            "84B6EB",
		"area-polaris":            "84B6EB",
		"exp-beginner":            "FEC5A8",
		"exp-expert":              "FEC5A8",
		"exp-intermediate":        "FEC5A8",
		"first-timer":             "50EF63",
		"help-wanted":             "169817",
		"kind-bug":                "B5142D",
		"kind-docs":               "B5142D",
		"kind-enhancement":        "B5142D",
		"kind-experimental":       "B5142D",
		"kind-feature":            "B5142D",
		"kind-question":           "B5142D",
		"platform-linux":          "CCCCCC",
		"platform-mac":            "CCCCCC",
		"platform-windows":        "CCCCCC",
		"priority-p0":             "FBC908",
		"priority-p1":             "FBC908",
		"priority-p2":             "FBC908",
		"priority-p3":             "FBC908",
		"resolution-by-design":    "FFDDDD",
		"resolution-duplicate":    "FFDDDD",
		"resolution-external":     "FFDDDD",
		"resolution-fixed":        "FFDDDD",
		"resolution-wont-fix":     "FFDDDD",
		"roadmap":                 "0B008C",
		"status-0-triage":         "004B75",
		"status-1-design-review":  "004B75",
		"status-2-code-review":    "004B75",
		"status-3-docs-review":    "004B75",
		"status-4-merge":          "004B75",
		"status-failing-ci":       "006C75",
		"status-more-info-needed": "006C75",
		"status-needs-attention":  "006C75",
		"status-needs-vendoring":  "006C75",
		"version-master":          "B2482A",
	}

	// Repository sort options -- Per Page 30 = default of GitHub UI
	repoOpt := &github.RepositoryListByOrgOptions{Type: "all", ListOptions: github.ListOptions{PerPage: 30}}

	githubOrganization := "rubrikinc"

	// Create an array to store the pagination results
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, githubOrganization, repoOpt)
		if err != nil {
			log.Fatal(err)
		}

		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}

		repoOpt.Page = resp.NextPage

	}

	for _, repo := range allRepos {
		repoName := repo.Name

		labels, _, err := client.Issues.ListLabels(ctx, githubOrganization, *repoName, &github.ListOptions{PerPage: 100})
		if err != nil {
			log.Fatal(err)
		}
		repoLabels := map[string]string{}
		for _, label := range labels {
			name := label.Name
			color := label.Color
			repoLabels[*name] = *color
		}

		if len(repoLabels) != len(buildLabels) {

			for name, color := range repoLabels {

				if deleteLabel(buildLabels, name, color) == true {
					_, err := client.Issues.DeleteLabel(ctx, githubOrganization, *repoName, name)
					if err != nil {
						log.Fatal(err)
					}
				}

			}

			for name, color := range buildLabels {

				if createLabel(repoLabels, name, color) == true {

					newLabel := github.Label{
						Name:  &name,
						Color: &color,
					}

					_, _, err := client.Issues.CreateLabel(ctx, githubOrganization, *repoName, &newLabel)
					if err != nil {
						log.Fatal(err)
					}
				}

			}

		}

	}

}
