package repos

import (
	"context"
	"fmt"
	"gitmic/api"
	"gitmic/internal/env"
	"gitmic/workerpool"
	"time"
)

func RunPool() error {
	ctx := context.TODO()

	wp := workerpool.NewPool(4, 4, 5 /*time.Second*1,*/, time.Millisecond*300)
	wp.RunBackground()

	ga, err := api.NewGitApi(wp, api.WithGitToken(env.GIT_TOKEN))
	if err != nil {
		return fmt.Errorf("new git api: %w", err)
	}

	if err := ga.GetReposByOrgPool(ctx, org, 100); err != nil {
		return fmt.Errorf("get repos on pool: %w", err)
	}

	return nil
}
