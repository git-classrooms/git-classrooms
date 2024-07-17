package worker

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockWorker struct {
	BaseWorker
	doWorkCalled int
}

func (mw *MockWorker) doWork() {
	mw.doWorkCalled++
}

func TestBaseWorker(t *testing.T) {
	t.Run("calls doWork periodically", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockWorker := &MockWorker{}
		mockWorker.BaseWorker.Worker = mockWorker

		workInterval := 10 * time.Millisecond
		mockWorker.Start(ctx, workInterval)

		// Warte einige Intervalle und überprüfe, ob doWork aufgerufen wurde
		time.Sleep(50 * time.Millisecond)

		assert.NotEqual(t, 0, mockWorker.doWorkCalled, "Expected doWork to be called at least once")
	})

	t.Run("stops calling doWork after context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockWorker := &MockWorker{}
		mockWorker.BaseWorker.Worker = mockWorker

		workInterval := 10 * time.Millisecond
		mockWorker.Start(ctx, workInterval)

		time.Sleep(50 * time.Millisecond)
		assert.Greater(t, mockWorker.doWorkCalled, 0, "The worker hasn't been called at least once")

		cancel()
		time.Sleep(20 * time.Millisecond) // let the worker stop

		callsBefore := mockWorker.doWorkCalled
		time.Sleep(50 * time.Millisecond) // wait a time and check that doWork has not been called again

		assert.Equal(t, callsBefore, mockWorker.doWorkCalled, "Worker has called doWork again, although context has been canceled")
	})
}
