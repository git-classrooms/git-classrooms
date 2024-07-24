package worker

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockWork struct {
	doCalled int
}

func (mw *MockWork) Do(tx context.Context) {
	mw.doCalled++
}

func TestWorker(t *testing.T) {
	t.Run("calls Do() periodically", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockWork := &MockWork{}
		worker := NewWorker(mockWork)

		workInterval := 10 * time.Millisecond
		worker.Start(ctx, workInterval)

		time.Sleep(50 * time.Millisecond) // wait for some intervals, if Do() has been called

		assert.NotEqual(t, 0, mockWork.doCalled, "Expected Do() to be called at least once")
	})

	t.Run("stops calling Do() after context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		mockWork := &MockWork{}
		worker := NewWorker(mockWork)

		workInterval := 10 * time.Millisecond
		worker.Start(ctx, workInterval)

		time.Sleep(50 * time.Millisecond)
		assert.Greater(t, mockWork.doCalled, 0, "The worker hasn't been called at least once")

		cancel()
		time.Sleep(20 * time.Millisecond) // let the worker stop

		callsBefore := mockWork.doCalled
		time.Sleep(50 * time.Millisecond) // wait a time and check that Do() has not been called again

		assert.Equal(t, callsBefore, mockWork.doCalled, "Worker has called Do() again, although context has been canceled")
	})
}
