package model

type TestReportTestCase struct {
	Status        string      `json:"status"`
	Name          string      `json:"name"`
	Classname     string      `json:"classname"`
	ExecutionTime float64     `json:"execution_time"`
	SystemOutput  interface{} `json:"system_output,omitempty"`
	StackTrace    string      `json:"stack_trace,omitempty"`
}
