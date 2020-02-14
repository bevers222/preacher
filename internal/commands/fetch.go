package commands

import (
	"context"
	"fmt"
	"github.com/google/go-github/v29/github"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// func fetch creates a list of all the repositories in the org, pulls the given file from the desired branch, and creates a folder and file for each repository
func Fetch(c *cli.Context) error {
	var numDownloaded uint32
	var numSkipped uint32
	var ctx = context.Background()
	var fileName = "Jenkinsfile"
	var orgName = c.String("org")
	var numRepos = 0
	var down = make(chan uint32, 1)
	var skip = make(chan uint32, 1)
	var path = c.String("directory")
	var debug = c.Bool("debug")
	var skippedRepos = make(chan string, 10000)
	if c.Args().First() != "" {
		fileName = c.Args().First()
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: c.String("token")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	down <- numDownloaded
	skip <- numSkipped
	var wg sync.WaitGroup
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, orgName, opt)
		if err != nil {
			return err
		}
		numRepos = numRepos + len(repos)
		for _, repo := range repos {
			repoName := repo.GetName()
			branchName := repo.GetDefaultBranch()
			if c.String("branch") != "" {
				branchName = c.String("branch")
			}
			time.Sleep(5 * time.Millisecond)
			wg.Add(1)
			go func(wg *sync.WaitGroup, repoName string, branchName string) {
				defer wg.Done()
				rawContent, _, _, err := client.Repositories.GetContents(ctx, orgName, repoName, fileName, &github.RepositoryContentGetOptions{Ref: branchName})
				if err != nil {
					if debug {
						fmt.Fprintf(c.App.Writer, "Error getting file from %v: %v\n", repoName, err)
					}
					select {
					case current := <-skip:
						skip <- (current + 1)
					case skippedRepos <- repoName:
					}
					return
				}
				content, err := rawContent.GetContent()
				if err != nil {
					if debug {
						fmt.Fprintf(c.App.Writer, "Error getting content from %v/%v: %v\n", repoName, fileName, err)
					}
					select {
					case current := <-skip:
						skip <- (current + 1)
					case skippedRepos <- repoName:
					}
					return
				}
				dirLoc := filepath.Join(path, repoName)
				fileLoc := filepath.Join(path, repoName, fileName)
				err = os.MkdirAll(dirLoc, os.ModePerm)
				if err != nil {
					if debug {
						fmt.Fprintf(c.App.Writer, "Error making directory: %v\n", dirLoc)
					}
					select {
					case current := <-skip:
						skip <- (current + 1)
					case skippedRepos <- repoName:
					}
					return
				}
				err = ioutil.WriteFile(fileLoc, []byte(content), 0644)
				if err != nil {
					if debug {
						fmt.Fprintf(c.App.Writer, "Error writing file: %v\n", fileLoc)
					}
					select {
					case current := <-skip:
						skip <- (current + 1)
					case skippedRepos <- repoName:
					}
					return
				}
				if debug {
					fmt.Fprintf(c.App.Writer, "File delivered: %v\n", fileLoc)
				}
				select {
				case current := <-down:
					down <- (current + 1)
				}
			}(&wg, repoName, branchName)
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	wg.Wait()
	close(skippedRepos)
	if debug {
		fmt.Fprintf(c.App.Writer, "\n\n\nSkipped Repositories:\n")
		for r := range skippedRepos {
			fmt.Fprintf(c.App.Writer, "%v\n", r)
		}
	}
	fmt.Fprintf(c.App.Writer, "\n\n\n")
	fmt.Fprintf(c.App.Writer, "total number of repositories found %v\n", numRepos)
	fmt.Fprintf(c.App.Writer, "total number of files delivered:  %v\n", <-down)
	fmt.Fprintf(c.App.Writer, "total number of repositories skipped:  %v\n", <-skip)
	return nil
}
