package api

import (
	"context"
	"fmt"
	"gitmic/api/orgs/repos"
	"gitmic/workerpool"
	"log"
	"net/http"
	"sync"

	"golang.org/x/sync/errgroup"
)

func (ga *GitApi) GetReposByOrgPool(ctx context.Context, org string, to int) error {
	wg := sync.WaitGroup{}

	for i := 1; i <= to; i++ {
		wg.Add(1)
		ii := i

		t := workerpool.NewTask(func() (interface{}, error) {
			defer wg.Done()

			repos, _, err := getReposByOrg(ctx, ga, org, ii, 1)
			if err != nil {
				log.Printf("err repos by org: %w", err)
				return nil, fmt.Errorf("get repos by org: %w", err)
			}

			fmt.Printf("%s\n", (*repos)[0].FullName)

			// time.Sleep(time.Second * 10)

			return repos, nil
		}, nil)

		ga.wp.AddTask(t)
	}

	wg.Wait()

	return nil
}

// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#list-organization-repositories
func (ga *GitApi) GetReposByOrg(ctx context.Context, org string, toPage, perPage int) (*[]*repos.Repo, error) {
	reposHundred, headers, err := getReposByOrg(ctx, ga, org, 1, perPage)
	if err != nil {
		return nil, fmt.Errorf("get repos from `%s` for first page by `%d`: %w", org, perPage, err)
	}

	pages, err := getAllPagesFromHeader(headers)
	if err != nil {
		return nil, fmt.Errorf("get all pages from header: %w", err)
	}

	if pages == 0 {
		return reposHundred, nil
	}

	dstPage := toPage
	if dstPage == 0 || dstPage > pages {
		dstPage = pages
	}

	cntRepos := dstPage * perPage
	reposAll := make([]*repos.Repo, cntRepos)

	copy(reposAll, *reposHundred)

	idxRepo := len(*reposHundred)
	g, ctx := errgroup.WithContext(ctx)
	mx := sync.Mutex{}

	for i := 2; i <= dstPage; i++ {
		ii := i

		g.Go(func() error {
			reposPage, _, err := getReposByOrg(ctx, ga, org, ii, perPage)
			if err != nil {
				return fmt.Errorf("get repos by org=`%s` page=%d perPage=%d: %w", org, ii, perPage, err)
			}

			mx.Lock()
			defer mx.Unlock()

			copy(reposAll[idxRepo:], *reposPage)
			idxRepo += len(*reposPage)

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("goroutine get repos: %w", err)
	}

	reposAll = reposAll[:idxRepo]

	return &reposAll, nil
}

func getReposByOrg(ctx context.Context, ga *GitApi, org string, curPage, perPage int) (*[]*repos.Repo, *http.Header, error) {
	req, err := repos.MakeRequest(ctx, ga.Host, org, curPage, perPage)
	if err != nil {
		return nil, nil, fmt.Errorf("make request: %w", err)
	}

	ga.prepareRequestToken(req)

	repos := make([]*repos.Repo, 0)

	headers, err := doRequest(req, &repos)
	if err != nil {
		return nil, nil, fmt.Errorf("do request: %w", err)
	}

	return &repos, headers, nil
}
