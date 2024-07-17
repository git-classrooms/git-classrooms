package worker

import (
	"context"
	"time"
)

type Worker interface {
	Start(ctx context.Context, workInterval time.Duration)
	doWork()
}

type BaseWorker struct {
	Worker
	ticker *time.Ticker
	ctx    context.Context
}

func (w *BaseWorker) Start(ctx context.Context, workInterval time.Duration) {
	w.ctx = ctx
	w.ticker = time.NewTicker(workInterval)

	go func() {
		for {
			select {
			case <-w.ctx.Done():
				w.ticker.Stop()
				return
			case <-w.ticker.C:
				w.doWork()
			}
		}
	}()
}
