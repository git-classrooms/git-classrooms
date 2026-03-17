package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
)

// ManualResult represents the result of a manual grading rubric.
type ManualResult struct {
	RubricID   uuid.UUID `json:"rubricId"`
	RubricName string    `json:"rubricName"`
	Score      int       `json:"score"`
	Feedback   string    `json:"feedback"`
	MaxScore   int       `json:"maxScore"`
}

// ReportDataItem represents a single item in a report.
type ReportDataItem struct {
	ProjectID           uuid.UUID               `json:"projectId"`
	AssignmentName      string                  `json:"assignmentName"`
	AssignmentDueDate   time.Time               `json:"assignmentDueDate"`
	TeamName            string                  `json:"teamName"`
	Name                string                  `json:"name"`
	Username            string                  `json:"username"`
	Email               string                  `json:"email"`
	RubricResults       map[string]ManualResult `json:"rubricResults"`
	AutogradingScore    int                     `json:"autogradingScore"`
	AutogradingMaxScore int                     `json:"autogradingMaxScore"`
	MaxScore            int                     `json:"maxScore"`
	Score               int                     `json:"score"`
	Percentage          float64                 `json:"percentage"`
}

// GenerateReports generates reports for the given assignments and rubrics.
func GenerateReports(assignments []*database.Assignment, rubrics []*database.ManualGradingRubric, teamID *uuid.UUID) ([][]*ReportDataItem, error) {
	reports := make([][]*ReportDataItem, len(assignments))

	for i, assignment := range assignments {
		report, err := GenerateReport(assignment, teamID)
		if err != nil {
			return nil, err
		}
		reports[i] = report
	}

	return reports, nil
}

// GenerateReport generates a report for the given assignment and rubrics.
func GenerateReport(assignment *database.Assignment, teamID *uuid.UUID) ([]*ReportDataItem, error) {
	reportData := createReportDataItems(assignment, teamID)

	return reportData, nil
}

// GenerateCSVReports generates CSV reports for the given assignments and rubrics.
func GenerateCSVReports(w io.Writer, assignments []*database.Assignment, rubrics []*database.ManualGradingRubric, teamID *uuid.UUID) error {
	writer := csv.NewWriter(w)
	writer.Comma = ';'
	err := writeHeader(writer, rubrics)
	if err != nil {
		return err
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}

	for _, assignment := range assignments {
		if err := GenerateCSVReport(w, assignment, rubrics, teamID, false); err != nil {
			return err
		}
	}

	return nil
}

// GenerateCSVReport generates a CSV report for the given assignment and rubrics.
func GenerateCSVReport(w io.Writer, assignment *database.Assignment, rubrics []*database.ManualGradingRubric, teamID *uuid.UUID, includeHeader bool) error {
	reportData := createReportDataItems(assignment, teamID)

	writer := csv.NewWriter(w)
	writer.Comma = ';'

	if includeHeader {
		err := writeHeader(writer, rubrics)
		if err != nil {
			return err
		}
	}

	for _, item := range reportData {
		row := []string{
			item.AssignmentName, item.TeamName, item.Name, item.Username, item.Email,
		}

		// Add manual rubric scores
		for _, rubric := range rubrics {
			result, ok := item.RubricResults[rubric.Name]
			if !ok {
				row = append(row, "", "", "")
				continue
			}

			row = append(row, strconv.Itoa(result.Score), result.Feedback, strconv.Itoa(result.MaxScore))
		}

		row = append(row, strconv.Itoa(item.AutogradingScore), strconv.Itoa(item.AutogradingMaxScore), strconv.Itoa(item.MaxScore), strconv.Itoa(item.Score), fmt.Sprintf("%.2f", item.Percentage))
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}
	return nil
}

// GenerateCSVReportForTeam generates a CSV report for the given assignment and rubrics for a specific team.
func createReportDataItems(assignment *database.Assignment, teamID *uuid.UUID) []*ReportDataItem {
	reportData := make([]*ReportDataItem, 0)

	for _, project := range assignment.Projects {
		if teamID != nil && project.TeamID != *teamID {
			continue
		}

		for _, projectGradingDate := range project.Gradings {
			assignmentDateIdx := slices.IndexFunc(assignment.AssignmentDates, func(a *database.AssignmentDate) bool { return projectGradingDate.AssignmentDateID == a.ID })
			if assignmentDateIdx < 0 {
				continue
			}

			assignmentDate := assignment.AssignmentDates[assignmentDateIdx]

			manualRubricResults := createManualRubricResults(projectGradingDate)
			autogradingScore := calculateAutogradingScore(projectGradingDate, assignmentDate.JUnitTests)
			autogradingMaxScore := calculateAutogradingMaxScore(projectGradingDate, assignmentDate.JUnitTests)

			maxScore := calculateMaxScore(projectGradingDate, assignmentDate.JUnitTests, assignmentDate.GradingManualRubrics)
			score := calculateScore(projectGradingDate, assignmentDate.JUnitTests)

			var percentage float64

			if maxScore == 0.0 {
				percentage = 0.0
			} else {
				percentage = float64(score) / float64(maxScore) * 100
			}

			for _, member := range project.Team.Member {
				reportData = append(reportData, &ReportDataItem{
					ProjectID:           project.ID,
					AssignmentName:      assignment.Name,
					AssignmentDueDate:   assignmentDate.DueDate,
					TeamName:            project.Team.Name,
					Name:                member.User.Name,
					Username:            member.User.GitlabUsername,
					Email:               member.User.GitlabEmail,
					RubricResults:       manualRubricResults,
					AutogradingScore:    autogradingScore,
					AutogradingMaxScore: autogradingMaxScore,
					MaxScore:            maxScore,
					Score:               score,
					Percentage:          percentage,
				})
			}
		}

	}

	slices.SortFunc(reportData, func(a, b *ReportDataItem) int {
		// Sort by assignment TeamName
		if a.TeamName != b.TeamName {
			return strings.Compare(a.TeamName, b.TeamName)
		}

		// Sort by assignment Member Name
		if a.Name != b.Name {
			return strings.Compare(a.Name, b.Name)
		}

		// Sort by assignment AssignmentDueDate
		return a.AssignmentDueDate.Compare(b.AssignmentDueDate)
	})

	return reportData
}

// createManualRubricResults creates a map of manual rubric results for a project.
func createManualRubricResults(project *database.AssignmentProjectGradingDate) map[string]ManualResult {
	results := make(map[string]ManualResult)
	for _, result := range project.GradingManualResults {
		feedback := ""
		if result.Feedback != nil {
			feedback = *result.Feedback
		}
		results[result.Rubric.Name] = ManualResult{
			RubricID:   result.RubricID,
			RubricName: result.Rubric.Name,
			Score:      result.Score,
			Feedback:   feedback,
			MaxScore:   result.Rubric.MaxScore,
		}
	}

	return results
}

// calculateMaxScore calculates the maximum score for a project.
func calculateMaxScore(project *database.AssignmentProjectGradingDate, tests []*database.AssignmentJunitTest, rubrics []*database.ManualGradingRubric) int {
	maxScore := 0
	for _, rubric := range rubrics {
		maxScore += rubric.MaxScore
	}
	return maxScore + calculateAutogradingMaxScore(project, tests)
}

// calculateAutogradingMaxScore calculates the maximum score for the autograding tests.
func calculateAutogradingMaxScore(project *database.AssignmentProjectGradingDate, tests []*database.AssignmentJunitTest) int {
	if len(tests) == 0 {
		if project.GradingJUnitTestResult != nil {
			return project.GradingJUnitTestResult.TotalCount
		}
		return 0
	}
	maxScore := 0
	for _, test := range tests {
		maxScore += test.Score
	}
	return maxScore
}

// calculateAutogradingScore calculates the score for the autograding tests.
func calculateAutogradingScore(project *database.AssignmentProjectGradingDate, tests []*database.AssignmentJunitTest) int {
	score := 0
	if project.GradingJUnitTestResult != nil {
		if len(tests) == 0 {
			return project.GradingJUnitTestResult.SuccessCount
		}

		for _, ts := range project.GradingJUnitTestResult.TestSuites {
			for _, tc := range ts.TestCases {
				name := fmt.Sprintf("%s/%s", ts.Name, tc.Name)
				index := slices.IndexFunc(tests, func(test *database.AssignmentJunitTest) bool {
					return test.Name == name
				})

				if index != -1 {
					if tc.Status == "success" {
						score += tests[index].Score
					}
				}
			}
		}
	}
	return score
}

// calculateScore calculates the score for a project.
func calculateScore(project *database.AssignmentProjectGradingDate, tests []*database.AssignmentJunitTest) int {
	score := 0
	for _, result := range project.GradingManualResults {
		score += result.Score
	}
	return score + calculateAutogradingScore(project, tests)
}

// writeHeader writes the header for the CSV report.
func writeHeader(writer *csv.Writer, rubrics []*database.ManualGradingRubric) error {
	header := []string{
		// 1               2	       3       4           5
		"AssignmentName", "TeamName", "Name", "Username", "Email",
	}

	// Add headers for manual rubric scores
	for _, rubric := range rubrics {
		header = append(header, rubric.Name+"Score", rubric.Name+"Feedback", rubric.Name+"MaxScore")
	}

	header = append(header, "AutogradingScore", "AutogradingMaxScore", "MaxScore", "Score", "Percentage")

	return writer.Write(header)
}
