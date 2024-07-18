package worker

import (
	"context"
	"time"
)

type Worker interface {
	Start(ctx context.Context, workInterval time.Duration)
	doWork(ctx context.Context)
}

type BaseWorker struct {
	Worker
	ticker *time.Ticker
}

func (w *BaseWorker) Start(ctx context.Context, workInterval time.Duration) {
	w.ticker = time.NewTicker(workInterval)

	go func() {
		for {
			select {
			case <-ctx.Done():
				w.ticker.Stop()
				return
			case <-w.ticker.C:
				w.doWork(ctx)
			}
		}
	}()
}
