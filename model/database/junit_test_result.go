package database

import "github.com/google/uuid"

type JUnitTestResult struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;not null;default:uuid_generate_v4()" json:"id"`
	TotalTime    float64   `json:"totalTime"`
	TotalCount   int       `json:"totalCount"`
	SuccessCount int       `json:"successCount"`
	FailedCount  int       `json:"failedCount"`
	SkippedCount int       `json:"skippedCount"`
	ErrorCount   int       `json:"errorCount"`
} //@Name JUnitTestResult
