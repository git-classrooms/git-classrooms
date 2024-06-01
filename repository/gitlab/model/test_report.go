package model

type TestReport struct {
	TotalTime    float64               `json:"total_time"`
	TotalCount   int                   `json:"total_count"`
	SuccessCount int                   `json:"success_count"`
	FailedCount  int                   `json:"failed_count"`
	SkippedCount int                   `json:"skipped_count"`
	ErrorCount   int                   `json:"error_count"`
	TestSuites   []TestReportTestSuite `json:"test_suites"`
}
