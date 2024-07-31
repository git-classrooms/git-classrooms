package utils

import (
	"encoding/csv"
	"fmt"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"slices"
	"strconv"
	"strings"
)

type ReportDataItem struct {
	AssignmentName   string            `json:"assignmentName,omitempty"`
	TeamName         string            `json:"teamName,omitempty"`
	Name             string            `json:"name,omitempty"`
	Username         string            `json:"username,omitempty"`
	Email            string            `json:"email,omitempty"`
	RubricScores     map[string]int    `json:"rubricScores,omitempty"`
	RubricFeedback   map[string]string `json:"rubricFeedback,omitempty"`
	AutogradingScore int               `json:"autogradingScore,omitempty"`
	MaxScore         int               `json:"maxScore,omitempty"`
	Score            int               `json:"score,omitempty"`
	Percentage       float64           `json:"percentage,omitempty"`
}

func GenerateReports(assignments []*database.Assignment) ([][]*ReportDataItem, error) {
	reports := make([][]*ReportDataItem, len(assignments))

	for i, assignment := range assignments {
		report, err := GenerateReport(assignment)
		if err != nil {
			return nil, err
		}
		reports[i] = report
	}

	return reports, nil
}

func GenerateReport(assignment *database.Assignment) ([]*ReportDataItem, error) {
	reportData := createReportDataItems(assignment)

	return reportData, nil
}

func GenerateCSVReports(assignments []*database.Assignment) ([]string, error) {
	reports := make([]string, len(assignments))

	for i, assignment := range assignments {
		report, err := GenerateCSVReport(assignment)
		if err != nil {
			return nil, err
		}
		reports[i] = report
	}

	return reports, nil
}

func GenerateCSVReport(assignment *database.Assignment) (string, error) {
	reportData := createReportDataItems(assignment)

	var b strings.Builder
	writer := csv.NewWriter(&b)

	header := []string{
		// 1               2	       3       4           5
		"AssignmentName", "TeamName", "Name", "Username", "Email",
	}

	// Add headers for manual rubric scores
	for _, rubric := range assignment.GradingManualRubrics {
		header = append(header, rubric.Name+"Score", rubric.Name+"Feedback")
	}

	header = append(header, "AutogradingScore", "MaxScore", "Score", "Percentage")

	err := writer.Write(header)
	if err != nil {
		return "", err
	}

	for _, item := range reportData {
		row := []string{
			item.AssignmentName, item.TeamName, item.Name, item.Username, item.Email,
		}

		// Add manual rubric scores
		for _, rubric := range assignment.GradingManualRubrics {
			row = append(row, strconv.Itoa(item.RubricScores[rubric.Name]), item.RubricFeedback[rubric.Name])
		}

		row = append(row, strconv.Itoa(item.AutogradingScore), strconv.Itoa(item.MaxScore), strconv.Itoa(item.Score), fmt.Sprintf("%.2f", item.Percentage))
		err := writer.Write(row)
		if err != nil {
			return "", err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}

	return b.String(), nil
}

func createReportDataItems(assignment *database.Assignment) []*ReportDataItem {
	reportData := make([]*ReportDataItem, 0)

	for _, project := range assignment.Projects {
		manualRubricScores := createManualRubricScoresMap(project)
		manualRubricFeedbacks := createManualRubricFeedbacksMap(project)
		autogradingScore := calculateAutogradingScore(project)
		maxScore := calculateMaxScore(project)
		score := calculateScore(project)
		percentage := float64(score) / float64(maxScore) * 100

		for _, member := range project.Team.Member {
			reportData = append(reportData, &ReportDataItem{
				AssignmentName:   project.Assignment.Name,
				TeamName:         project.Team.Name,
				Name:             member.User.Name,
				Username:         member.User.GitlabUsername,
				Email:            member.User.GitlabEmail,
				RubricScores:     manualRubricScores,
				RubricFeedback:   manualRubricFeedbacks,
				AutogradingScore: autogradingScore,
				MaxScore:         maxScore,
				Score:            score,
				Percentage:       percentage,
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

func calculateMaxScore(project *database.AssignmentProjects) int {
	maxScore := 0
	for _, result := range project.GradingManualResults {
		maxScore += result.Rubric.MaxScore
	}
	if project.GradingJUnitTestResult != nil {
		maxScore += project.GradingJUnitTestResult.TotalCount
	}
	return maxScore
}

func calculateAutogradingScore(project *database.AssignmentProjects) int {
	if project.GradingJUnitTestResult != nil {
		return project.GradingJUnitTestResult.SuccessCount
	}
	return 0
}

func calculateScore(project *database.AssignmentProjects) int {
	score := 0
	for _, result := range project.GradingManualResults {
		score += result.Score
	}
	if project.GradingJUnitTestResult != nil {
		score += project.GradingJUnitTestResult.SuccessCount
	}
	return score
}
