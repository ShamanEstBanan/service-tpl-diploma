package workers

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"service-tpl-diploma/internal/domain"
)

func RunPool(ctx context.Context, size int, jobs chan domain.Job) error {
	gr, ctx := errgroup.WithContext(ctx)
	for i := 0; i < size; i++ {
		gr.Go(func() error {
			for {
				select {
				case job := <-jobs:
					err := job.Run(ctx)
					if err != nil {
						fmt.Printf("Job error: %s \n", err)
						return err
					}
				case <-ctx.Done():
					fmt.Println("Context canceled")
					return ctx.Err()
				}
			}
		})
	}

	return gr.Wait()
}
