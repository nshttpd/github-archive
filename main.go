package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"golang.org/x/tools/go/vcs"
)

type config struct {
	oauthToken string
	ghOrg      string
	ghRepo     string
	ghUser     string
	maxFunc    int
	ts         oauth2.TokenSource
}

func archive(r *github.Repository, c *config, wg *sync.WaitGroup, t chan int) {

	if wg != nil {
		defer wg.Done()
	}

	path := fmt.Sprintf("./%s", *r.Name)
	repo := fmt.Sprintf("%s@github.com:%s/%s", c.ghUser, c.ghOrg, *r.Name)

	cmd := vcs.ByCmd("git")

	// if it doesn't exist go and clone it
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("\tcloning %s\n", repo)
		err := cmd.Create(path, repo)
		if err != nil {
			fmt.Printf("error on archive call for repository : %s\n", *r.Name)
			fmt.Printf("error : %s\n", err)

		}

	} else { // if it does exist do a pull
		fmt.Printf("\tupdating %s\n", repo)
		cmd.Download(path)
	}

	if t != nil {
		<-t
	}

}

func getRepos(c *config) ([]*github.Repository, error) {

	ctx := context.Background()
	tc := oauth2.NewClient(ctx, c.ts)
	client := github.NewClient(tc)

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 20},
	}

	var allRepos []*github.Repository
	fmt.Print("fetching repos.")
	for {

		repos, resp, err := client.Repositories.ListByOrg(context.Background(), c.ghOrg, opt)

		if err != nil {
			return nil, err
		}

		allRepos = append(allRepos, repos...)

		if resp.NextPage == 0 {
			fmt.Println(" done")
			break
		}

		opt.ListOptions.Page = resp.NextPage
		fmt.Print(".")

	}

	return allRepos, nil
}

func main() {
	cfg := &config{}

	flag.StringVar(&cfg.oauthToken, "token", "", "oauth token for api access")
	flag.StringVar(&cfg.ghOrg, "org", "", "organization to list")
	flag.StringVar(&cfg.ghRepo, "repo", "", "single repo to archive if not all")
	flag.StringVar(&cfg.ghUser, "user", "git", "user to connect to github as via ssh")
	flag.IntVar(&cfg.maxFunc, "max", 4, "max goprocs to fetch with")
	flag.Parse()

	if cfg.oauthToken == "" {
		fmt.Println("must supply oauth token to authenticate with")
		os.Exit(1)
	}

	if cfg.ghOrg == "" {
		fmt.Println("must specify an organization")
		os.Exit(1)
	}

	cfg.ts = oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.oauthToken},
	)

	if cfg.ghRepo != "" {
		r := &github.Repository{
			Name: &cfg.ghRepo,
		}
		archive(r, cfg, nil, nil)
	} else {
		repos, err := getRepos(cfg)
		if err != nil {
			fmt.Println("error fetching list of repos")
			fmt.Printf("error: %s\n", err)
			os.Exit(1)
		}

		throttle := make(chan int, cfg.maxFunc)
		var wg sync.WaitGroup

		for _, r := range repos {
			throttle <- 1
			wg.Add(1)
			fmt.Printf("repo : %s\n", *r.Name)
			go archive(r, cfg, &wg, throttle)
		}
		wg.Wait()
	}

	os.Exit(0)
}
