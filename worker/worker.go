package worker

import (
	"context"
	"time"
)

type Work interface {
	Do(context.Context)
}

type Worker struct {
	work Work
}

func NewWorker(work Work) *Worker {
	return &Worker{work}
}

func (w *Worker) Start(ctx context.Context, workInterval time.Duration) {
	ticker := time.NewTicker(1 * time.Millisecond)
	first := true

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if first {
					ticker.Reset(workInterval)
				}
				w.work.Do(ctx)
			}
		}
	}()
}
