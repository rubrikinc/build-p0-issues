package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {

	ctx := context.Background()

	githubToken := os.Getenv("GITHUB_AUTH")
	// Authentication configuration
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// Repository sort options -- Per Page 30 = default of GitHub UI
	repoOpt := &github.RepositoryListByOrgOptions{Type: "public", ListOptions: github.ListOptions{PerPage: 30}}

	// Issue Options
	issueOpt := &github.IssueListByRepoOptions{Labels: []string{"priority-p0"}, State: "open"}

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

		issues, _, err := client.Issues.ListByRepo(ctx, "rubrikinc", *repoName, issueOpt)
		if err != nil {
			log.Fatal(err)
		}

		if len(issues) != 0 {
			fmt.Printf("\n%s\n", *repoName)

		}

		for _, issue := range issues {
			issueTitle := issue.Title
			fmt.Printf("- %s\n", *issueTitle)

		}

	}

}
