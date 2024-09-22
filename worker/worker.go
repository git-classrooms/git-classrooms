// Package worker provides functionality for automating tasks.
//
// The main components of the package include:
// - DueAssignmentWork: Handles the closure of assignments that have passed their due date.
// - SyncGitlabDbWork: Synchronizes classrooms, teams, and projects between the local database and GitLab.
// - Worker: Provides a mechanism to run tasks periodically at specified intervals.
package worker

import (
	"context"
	"time"
)

// Work interface defines a single method, Do(), which performs a task using the provided context.
type Work interface {
	Do(context.Context)
}

// Worker is responsible for executing a piece of work periodically.
type Worker struct {
	work Work
}

// NewWorker creates a new Worker instance with the provided work to be done.
func NewWorker(work Work) *Worker {
	return &Worker{work}
}

// Start begins the periodic execution of the work at the specified interval.
// It runs until the provided context is canceled.
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
