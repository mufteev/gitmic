package repos

import (
	"context"
	"fmt"
	"gitmic/api"
	"gitmic/internal/env"
	"gitmic/semaphore"
	"gitmic/workerpool"
	"time"
)

func RunPool() error {
	ctx := context.TODO()

	sem := semaphore.NewSemaphore(4, time.Millisecond*500, time.Millisecond*800)

	wp := workerpool.NewPool(4, 4, sem)
	wp.RunBackground()

	ga, err := api.NewGitApi(wp, api.WithGitToken(env.GIT_TOKEN))
	if err != nil {
		return fmt.Errorf("new git api: %w", err)
	}

	if err := ga.GetReposByOrgPool(ctx, org, 50); err != nil {
		return fmt.Errorf("get repos on pool: %w", err)
	}

	return nil
}
