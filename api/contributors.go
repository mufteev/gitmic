package api

import (
	"context"
	"fmt"
	"gitmic/api/repos/contributors"
	"net/http"
	"sync"

	"golang.org/x/sync/errgroup"
)

// https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#list-repository-contributors
func (ga *GitApi) GetContributorsByRepo(ctx context.Context, repo string, toPage, perPage int) (*[]*contributors.Contributor, error) {
	contribs, headers, err := getContributorsByRepo(ctx, ga, repo, 1, perPage)
	if err != nil {
		return nil, fmt.Errorf("get contribs from `%s` for first page by `%d`: %w", repo, perPage, err)
	}

	pages, err := getAllPagesFromHeader(headers)
	if err != nil {
		return nil, fmt.Errorf("get all pages from header: %w", err)
	}

	if pages == 0 {
		return contribs, nil
	}

	dstPage := toPage
	if dstPage == 0 || dstPage > pages {
		dstPage = pages
	}

	cntContribs := dstPage * perPage
	contribsAll := make([]*contributors.Contributor, cntContribs)

	copy(contribsAll, *contribs)

	idxContrib := len(*contribs)
	g, ctx := errgroup.WithContext(ctx)
	mx := sync.Mutex{}

	for i := 2; i <= dstPage; i++ {
		ii := i

		g.Go(func() error {
			contribPage, _, err := getContributorsByRepo(ctx, ga, repo, ii, perPage)
			if err != nil {
				return fmt.Errorf("get repos by reop=`%s` page=%d perPage=%d: %w", repo, ii, perPage, err)
			}

			mx.Lock()
			defer mx.Unlock()

			copy(contribsAll[idxContrib:], *contribPage)
			idxContrib += len(*contribPage)

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("goroutine get repos: %w", err)
	}

	contribsAll = contribsAll[:idxContrib]

	return &contribsAll, nil
}

func getContributorsByRepo(ctx context.Context, ga *GitApi, repo string, toPage, perPage int) (*[]*contributors.Contributor, *http.Header, error) {
	req, err := contributors.MakeRequest(ctx, ga.Host, repo, toPage, perPage)
	if err != nil {
		return nil, nil, fmt.Errorf("make request: %w", err)
	}

	ga.prepareRequestToken(req)

	contribs := make([]*contributors.Contributor, 0)

	headers, err := doRequest(req, &contribs)
	if err != nil {
		return nil, nil, fmt.Errorf("do request: %w", err)
	}

	return &contribs, headers, nil
}
