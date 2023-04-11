package worker

import (
	"context"

	"parser/internal/services/crawler"
)

func (w worker) Produce(ctx context.Context, crawlers []crawler.Crawler) {
	go func() {
		for i, c := range crawlers {
			err := w.rl.Wait(ctx)
			if err != nil {
				w.log.WithError(err).Fatal("failed to wait for the rate limiter")
			}

			w.tasksPending.Inc()
			w.tasksChan <- Task{
				seq: i,
				do:  c.Crawl,
			}
		}

		defer func() {
			close(w.tasksChan)
		}()
	}()
}
