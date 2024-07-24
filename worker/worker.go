package worker

import (
	"context"
	"time"
)

type Work interface {
	Do(context.Context)
}

type Worker struct {
	work   Work
	ticker *time.Ticker
}

func NewWorker(work Work) *Worker {
	return &Worker{work, nil}
}

func (w *Worker) Start(ctx context.Context, workInterval time.Duration) {
	w.ticker = time.NewTicker(workInterval)

	go func() {
		for {
			select {
			case <-ctx.Done():
				w.ticker.Stop()
				return
			case <-w.ticker.C:
				w.work.Do(ctx)
			}
		}
	}()
}
