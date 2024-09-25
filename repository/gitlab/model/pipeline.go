package model

import "time"

type Pipeline struct {
	ID          int
	Status      string
	Ref         string
	UpdatedAt   *time.Time
	CreatedAt   *time.Time
	StartedAt   *time.Time
	FinishedAt  *time.Time
	CommittedAt *time.Time
	Duration    int
	WebURL      string
}
