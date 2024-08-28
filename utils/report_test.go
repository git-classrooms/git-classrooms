package utils

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
)

func TestGenerateReports(t *testing.T) {
	assignments := []*database.Assignment{
		{
			Name: "Assignment 1",
			GradingManualRubrics: []*database.ManualGradingRubric{
				{Name: "Quality", MaxScore: 10},
			},
			Projects: []*database.AssignmentProjects{
				{
					Assignment: database.Assignment{Name: "Assignment 1"},
					GradingManualResults: []*database.ManualGradingResult{
						{Rubric: database.ManualGradingRubric{Name: "Quality", MaxScore: 10}, Score: 8},
					},
					GradingJUnitTestResult: &database.JUnitTestResult{TestReport: model.TestReport{TotalCount: 5, SuccessCount: 4}},
					Team: database.Team{
						Name: "Team A",
						Member: []*database.UserClassrooms{
							{
								User: database.User{Name: "John Doe", GitlabUsername: "johndoe", GitlabEmail: "john.doe@example.com"},
							},
						},
					},
				},
			},
		},
	}

	reports, err := GenerateReports(assignments, nil)
	assert.NoError(t, err)
	assert.NotNil(t, reports)
	assert.Equal(t, 1, len(reports))
	assert.Equal(t, "Team A", reports[0][0].TeamName)
	assert.Equal(t, "Assignment 1", reports[0][0].AssignmentName)
}

func TestGenerateCSVReports(t *testing.T) {
	assignments := []*database.Assignment{
		{
			Name: "Assignment 1",
			GradingManualRubrics: []*database.ManualGradingRubric{
				{Name: "Quality", MaxScore: 10},
			},
			Projects: []*database.AssignmentProjects{
				{
					Assignment: database.Assignment{Name: "Assignment 1"},
					GradingManualResults: []*database.ManualGradingResult{
						{Rubric: database.ManualGradingRubric{Name: "Quality", MaxScore: 10}, Score: 8},
					},
					GradingJUnitTestResult: &database.JUnitTestResult{TestReport: model.TestReport{TotalCount: 5, SuccessCount: 4}},
					Team: database.Team{
						Name: "Team A",
						Member: []*database.UserClassrooms{
							{
								User: database.User{Name: "John Doe", GitlabUsername: "johndoe", GitlabEmail: "john.doe@example.com"},
							},
						},
					},
				},
			},
		},
		{
			Name: "Assignment 2",
			GradingManualRubrics: []*database.ManualGradingRubric{
				{Name: "Beauty", MaxScore: 20},
			},
			Projects: []*database.AssignmentProjects{
				{
					Assignment: database.Assignment{Name: "Assignment 2"},
					GradingManualResults: []*database.ManualGradingResult{
						{Rubric: database.ManualGradingRubric{Name: "Beauty", MaxScore: 20}, Score: 15},
					},
					GradingJUnitTestResult: &database.JUnitTestResult{TestReport: model.TestReport{TotalCount: 11, SuccessCount: 7}},
					Team: database.Team{
						Name: "Team A",
						Member: []*database.UserClassrooms{
							{
								User: database.User{Name: "John Doe", GitlabUsername: "johndoe", GitlabEmail: "john.doe@example.com"},
							},
						},
					},
				},
			},
		},
	}

	var buffer bytes.Buffer
	err := GenerateCSVReports(&buffer, assignments, nil)
	assert.NoError(t, err)
	r, err := zip.NewReader(bytes.NewReader(buffer.Bytes()), int64(buffer.Len()))
	assert.NoError(t, err)
	assert.Equal(t, 2, len(r.File))

	csvRecords := make([][][]string, len(r.File))

	for i, file := range r.File {
		rc, err := file.Open()
		assert.NoError(t, err)

		csvReader := csv.NewReader(rc)
		records, err := csvReader.ReadAll()
		assert.NoError(t, err)
		csvRecords[i] = records
		rc.Close()
	}

	records := csvRecords[0]
	assert.Equal(t, "AssignmentName", records[0][0])
	assert.Equal(t, "Assignment 1", records[1][0])
	assert.Equal(t, "TeamName", records[0][1])
	assert.Equal(t, "Name", records[0][2])
	assert.Equal(t, "Username", records[0][3])
	assert.Equal(t, "Email", records[0][4])
	assert.Equal(t, "QualityScore", records[0][5])
	assert.Equal(t, "8", records[1][5])
	assert.Equal(t, "QualityFeedback", records[0][6])
	assert.Equal(t, "", records[1][6])
	assert.Equal(t, "AutogradingScore", records[0][7])
	assert.Equal(t, "4", records[1][7])
	assert.Equal(t, "MaxScore", records[0][8])
	assert.Equal(t, "15", records[1][8])
	assert.Equal(t, "Score", records[0][9])
	assert.Equal(t, "12", records[1][9])
	assert.Equal(t, "Percentage", records[0][10])

	records = csvRecords[1]
	assert.Equal(t, "AssignmentName", records[0][0])
	assert.Equal(t, "Assignment 2", records[1][0])
	assert.Equal(t, "TeamName", records[0][1])
	assert.Equal(t, "Name", records[0][2])
	assert.Equal(t, "Username", records[0][3])
	assert.Equal(t, "Email", records[0][4])
	assert.Equal(t, "BeautyScore", records[0][5])
	assert.Equal(t, "15", records[1][5])
	assert.Equal(t, "BeautyFeedback", records[0][6])
	assert.Equal(t, "", records[1][6])
	assert.Equal(t, "AutogradingScore", records[0][7])
	assert.Equal(t, "7", records[1][7])
	assert.Equal(t, "MaxScore", records[0][8])
	assert.Equal(t, "31", records[1][8])
	assert.Equal(t, "Score", records[0][9])
	assert.Equal(t, "22", records[1][9])
	assert.Equal(t, "Percentage", records[0][10])
}

func TestGenerateCSVReport(t *testing.T) {
	assignments := &database.Assignment{
		Name: "Assignment 1",
		GradingManualRubrics: []*database.ManualGradingRubric{
			{Name: "Quality", MaxScore: 10},
		},
		Projects: []*database.AssignmentProjects{
			{
				Assignment: database.Assignment{Name: "Assignment 1"},
				GradingManualResults: []*database.ManualGradingResult{
					{Rubric: database.ManualGradingRubric{Name: "Quality", MaxScore: 10}, Score: 8},
				},
				GradingJUnitTestResult: &database.JUnitTestResult{TestReport: model.TestReport{TotalCount: 5, SuccessCount: 4}},
				Team: database.Team{
					Name: "Team A",
					Member: []*database.UserClassrooms{
						{
							User: database.User{Name: "John Doe", GitlabUsername: "johndoe", GitlabEmail: "john.doe@example.com"},
						},
					},
				},
			},
		},
	}

	var reports strings.Builder

	err := GenerateCSVReport(&reports, assignments, nil)
	csvReports := reports.String()
	assert.NoError(t, err)
	assert.NotNil(t, csvReports)

	r := csv.NewReader(strings.NewReader(csvReports))
	records, err := r.ReadAll()
	assert.NoError(t, err)
	assert.Equal(t, "AssignmentName", records[0][0])
	assert.Equal(t, "Assignment 1", records[1][0])
	assert.Equal(t, "TeamName", records[0][1])
	assert.Equal(t, "Name", records[0][2])
	assert.Equal(t, "Username", records[0][3])
	assert.Equal(t, "Email", records[0][4])
	assert.Equal(t, "QualityScore", records[0][5])
	assert.Equal(t, "8", records[1][5])
	assert.Equal(t, "QualityFeedback", records[0][6])
	assert.Equal(t, "", records[1][6])
	assert.Equal(t, "AutogradingScore", records[0][7])
	assert.Equal(t, "4", records[1][7])
	assert.Equal(t, "MaxScore", records[0][8])
	assert.Equal(t, "15", records[1][8])
	assert.Equal(t, "Score", records[0][9])
	assert.Equal(t, "12", records[1][9])
	assert.Equal(t, "Percentage", records[0][10])
}
