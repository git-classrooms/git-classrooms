package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
)

type ReportDataItem struct {
	AssignmentName      string            `json:"assignmentName"`
	TeamName            string            `json:"teamName"`
	Name                string            `json:"name"`
	Username            string            `json:"username"`
	Email               string            `json:"email"`
	RubricScores        map[string]int    `json:"rubricScores"`
	RubricFeedback      map[string]string `json:"rubricFeedback"`
	RubricMaxScores     map[string]int    `json:"rubricMaxScores"`
	AutogradingScore    int               `json:"autogradingScore"`
	AutogradingMaxScore int               `json:"autogradingMaxScore"`
	MaxScore            int               `json:"maxScore"`
	Score               int               `json:"score"`
	Percentage          float64           `json:"percentage"`
}

func GenerateReports(assignments []*database.Assignment, rubrics []*database.ManualGradingRubric, teamID *uuid.UUID) ([][]*ReportDataItem, error) {
	reports := make([][]*ReportDataItem, len(assignments))

	for i, assignment := range assignments {
		report, err := GenerateReport(assignment, rubrics, teamID)
		if err != nil {
			return nil, err
		}
		reports[i] = report
	}

	return reports, nil
}

func GenerateReport(assignment *database.Assignment, rubrics []*database.ManualGradingRubric, teamID *uuid.UUID) ([]*ReportDataItem, error) {
	reportData := createReportDataItems(assignment, teamID)

	return reportData, nil
}

func GenerateCSVReports(w io.Writer, assignments []*database.Assignment, rubrics []*database.ManualGradingRubric, teamID *uuid.UUID) error {
	writer := csv.NewWriter(w)
	writer.Comma = ';'
	writeHeader(writer, rubrics)
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

func GenerateCSVReport(w io.Writer, assignment *database.Assignment, rubrics []*database.ManualGradingRubric, teamID *uuid.UUID, includeHeader bool) error {
	reportData := createReportDataItems(assignment, teamID)

	writer := csv.NewWriter(w)
	writer.Comma = ';'

	if includeHeader {
		writeHeader(writer, rubrics)
	}

	for _, item := range reportData {
		row := []string{
			item.AssignmentName, item.TeamName, item.Name, item.Username, item.Email,
		}

		// Add manual rubric scores
		for _, rubric := range rubrics {
			score, scoreOk := item.RubricScores[rubric.Name]
			feedback, feedbackOk := item.RubricFeedback[rubric.Name]
			maxScore, maxScoreOk := item.RubricMaxScores[rubric.Name]
			if !scoreOk || !feedbackOk || !maxScoreOk {
				row = append(row, "", "", "")
				continue
			}

			row = append(row, strconv.Itoa(score), feedback, strconv.Itoa(maxScore))
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

func createReportDataItems(assignment *database.Assignment, teamID *uuid.UUID) []*ReportDataItem {
	reportData := make([]*ReportDataItem, 0)

	for _, project := range assignment.Projects {
		if teamID != nil && project.TeamID != *teamID {
			continue
		}
		manualRubricScores := createManualRubricScoresMap(project)
		manualRubricFeedbacks := createManualRubricFeedbacksMap(project)
		manualRubricMaxScores := createManualRubricMaxScoresMap(project)
		autogradingScore := calculateAutogradingScore(project, assignment.JUnitTests)
		autogradingMaxScore := calculateAutogradingMaxScore(project, assignment.JUnitTests)
		maxScore := calculateMaxScore(project, assignment.JUnitTests)
		score := calculateScore(project, assignment.JUnitTests)
		var percentage float64

		if maxScore == 0.0 {
			percentage = 0.0
		} else {
			percentage = float64(score) / float64(maxScore) * 100
		}

		for _, member := range project.Team.Member {
			reportData = append(reportData, &ReportDataItem{
				AssignmentName:      assignment.Name,
				TeamName:            project.Team.Name,
				Name:                member.User.Name,
				Username:            member.User.GitlabUsername,
				Email:               member.User.GitlabEmail,
				RubricScores:        manualRubricScores,
				RubricFeedback:      manualRubricFeedbacks,
				RubricMaxScores:     manualRubricMaxScores,
				AutogradingScore:    autogradingScore,
				AutogradingMaxScore: autogradingMaxScore,
				MaxScore:            maxScore,
				Score:               score,
				Percentage:          percentage,
			})
		}
	}

	slices.SortFunc(reportData, func(a, b *ReportDataItem) int {
		// Sort by assignment TeamName
		if a.TeamName != b.TeamName {
			return strings.Compare(a.TeamName, b.TeamName)
		}

		// Sort by assignment Member Name
		return strings.Compare(a.Name, b.Name)
	})

	return reportData
}

func createManualRubricScoresMap(project *database.AssignmentProjects) map[string]int {
	points := make(map[string]int)
	for _, result := range project.GradingManualResults {
		points[result.Rubric.Name] = result.Score
	}
	return points
}
func createManualRubricMaxScoresMap(project *database.AssignmentProjects) map[string]int {
	points := make(map[string]int)
	for _, result := range project.GradingManualResults {
		points[result.Rubric.Name] = result.Rubric.MaxScore
	}
	return points
}

func createManualRubricFeedbacksMap(project *database.AssignmentProjects) map[string]string {
	feedback := make(map[string]string)
	for _, result := range project.GradingManualResults {
		if result.Feedback != nil {
			feedback[result.Rubric.Name] = *result.Feedback
		} else {
			feedback[result.Rubric.Name] = ""
		}
	}
	return feedback
}

func calculateMaxScore(project *database.AssignmentProjects, tests []*database.AssignmentJunitTest) int {
	maxScore := 0
	for _, result := range project.GradingManualResults {
		maxScore += result.Rubric.MaxScore
	}
	return maxScore + calculateAutogradingMaxScore(project, tests)
}

func calculateAutogradingMaxScore(project *database.AssignmentProjects, tests []*database.AssignmentJunitTest) int {
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

func calculateAutogradingScore(project *database.AssignmentProjects, tests []*database.AssignmentJunitTest) int {
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

func calculateScore(project *database.AssignmentProjects, tests []*database.AssignmentJunitTest) int {
	score := 0
	for _, result := range project.GradingManualResults {
		score += result.Score
	}
	return score + calculateAutogradingScore(project, tests)
}

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
