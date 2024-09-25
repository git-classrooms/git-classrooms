package utils

import (
	"bytes"
	"encoding/csv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
)

var (
	gradingJUnitTestResult = database.JUnitTestResult{TestReport: model.TestReport{TotalCount: 5, SuccessCount: 4, TestSuites: []model.TestReportTestSuite{{Name: "golang", TestCases: []model.TestReportTestCase{
		{Name: "test", Status: "success"},
		{Name: "test2", Status: "failed"},
		{Name: "test3", Status: "failed"},
		{Name: "test4", Status: "success"},
	}}}}}
	gradingJUnitTestResult2 = database.JUnitTestResult{TestReport: model.TestReport{TotalCount: 11, SuccessCount: 7, TestSuites: []model.TestReportTestSuite{{Name: "golang", TestCases: []model.TestReportTestCase{{Name: "test", Status: "success"}}}}}}
)

func TestGenerateReports(t *testing.T) {
	rubrics := []*database.ManualGradingRubric{
		{Name: "Quality", MaxScore: 10},
	}
	assignments := []*database.Assignment{
		{
			Name:                 "Assignment 1",
			GradingManualRubrics: rubrics,
			Projects: []*database.AssignmentProjects{
				{
					Assignment: database.Assignment{Name: "Assignment 1"},
					GradingManualResults: []*database.ManualGradingResult{
						{Rubric: *rubrics[0], Score: 8},
					},
					GradingJUnitTestResult: &gradingJUnitTestResult,
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

	reports, err := GenerateReports(assignments, rubrics, nil)
	assert.NoError(t, err)
	assert.NotNil(t, reports)
	assert.Equal(t, 1, len(reports))
	assert.Equal(t, "Team A", reports[0][0].TeamName)
	assert.Equal(t, "Assignment 1", reports[0][0].AssignmentName)
}

func TestGenerateCSVReports(t *testing.T) {
	rubrics := []*database.ManualGradingRubric{
		{Name: "Quality", MaxScore: 10},
	}

	rubrics2 := []*database.ManualGradingRubric{
		{Name: "Beauty", MaxScore: 20},
	}
	assignments := []*database.Assignment{
		{
			Name:                 "Assignment 1",
			GradingManualRubrics: rubrics,
			Projects: []*database.AssignmentProjects{
				{
					Assignment: database.Assignment{Name: "Assignment 1"},
					GradingManualResults: []*database.ManualGradingResult{
						{Rubric: *rubrics[0], Score: 8},
					},
					GradingJUnitTestResult: &gradingJUnitTestResult,
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
			Name:                 "Assignment 2",
			GradingManualRubrics: rubrics2,
			Projects: []*database.AssignmentProjects{
				{
					Assignment: database.Assignment{Name: "Assignment 2"},
					GradingManualResults: []*database.ManualGradingResult{
						{Rubric: *rubrics2[0], Score: 15},
					},
					GradingJUnitTestResult: &gradingJUnitTestResult2,
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

	classroomRubrics := make([]*database.ManualGradingRubric, 0)
	classroomRubrics = append(classroomRubrics, rubrics...)
	classroomRubrics = append(classroomRubrics, rubrics2...)

	var buffer bytes.Buffer
	err := GenerateCSVReports(&buffer, assignments, classroomRubrics, nil)
	assert.NoError(t, err)

	csvReader := csv.NewReader(bytes.NewReader(buffer.Bytes()))
	csvReader.Comma = ';'
	records, err := csvReader.ReadAll()
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
	assert.Equal(t, "QualityMaxScore", records[0][7])
	assert.Equal(t, "10", records[1][7])

	assert.Equal(t, "BeautyScore", records[0][8])
	assert.Equal(t, "", records[1][8])
	assert.Equal(t, "BeautyFeedback", records[0][9])
	assert.Equal(t, "", records[1][9])
	assert.Equal(t, "BeautyMaxScore", records[0][10])
	assert.Equal(t, "", records[1][10])

	assert.Equal(t, "AutogradingScore", records[0][11])
	assert.Equal(t, "4", records[1][11])
	assert.Equal(t, "AutogradingMaxScore", records[0][12])
	assert.Equal(t, "5", records[1][12])
	assert.Equal(t, "MaxScore", records[0][13])
	assert.Equal(t, "15", records[1][13])
	assert.Equal(t, "Score", records[0][14])
	assert.Equal(t, "12", records[1][14])
	assert.Equal(t, "Percentage", records[0][15])

	// --------------------------------------
	assert.Equal(t, "Assignment 2", records[2][0])

	assert.Equal(t, "", records[2][5])
	assert.Equal(t, "", records[2][6])
	assert.Equal(t, "", records[2][7])

	assert.Equal(t, "15", records[2][8])
	assert.Equal(t, "", records[2][9])
	assert.Equal(t, "20", records[2][10])

	assert.Equal(t, "7", records[2][11])
	assert.Equal(t, "11", records[2][12])
	assert.Equal(t, "31", records[2][13])
	assert.Equal(t, "22", records[2][14])
}

func TestGenerateCSVReport(t *testing.T) {
	rubrics := []*database.ManualGradingRubric{
		{Name: "Quality", MaxScore: 10},
	}
	assignments := &database.Assignment{
		Name:                 "Assignment 1",
		GradingManualRubrics: rubrics,
		Projects: []*database.AssignmentProjects{
			{
				Assignment: database.Assignment{Name: "Assignment 1"},
				GradingManualResults: []*database.ManualGradingResult{
					{Rubric: *rubrics[0], Score: 8},
				},
				GradingJUnitTestResult: &gradingJUnitTestResult,
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

	err := GenerateCSVReport(&reports, assignments, rubrics, nil, true)
	csvReports := reports.String()
	assert.NoError(t, err)
	assert.NotNil(t, csvReports)

	r := csv.NewReader(strings.NewReader(csvReports))
	r.Comma = ';'
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
	assert.Equal(t, "QualityMaxScore", records[0][7])
	assert.Equal(t, "10", records[1][7])
	assert.Equal(t, "AutogradingScore", records[0][8])
	assert.Equal(t, "4", records[1][8])
	assert.Equal(t, "AutogradingMaxScore", records[0][9])
	assert.Equal(t, "5", records[1][9])
	assert.Equal(t, "MaxScore", records[0][10])
	assert.Equal(t, "15", records[1][10])
	assert.Equal(t, "Score", records[0][11])
	assert.Equal(t, "12", records[1][11])
	assert.Equal(t, "Percentage", records[0][12])
}

func TestGenerateCSVReportWithTestScores(t *testing.T) {
	rubrics := []*database.ManualGradingRubric{
		{Name: "Quality", MaxScore: 10},
	}
	assignments := &database.Assignment{
		Name:                 "Assignment 1",
		JUnitTests:           []*database.AssignmentJunitTest{{Name: "golang/test", Score: 2}, {Name: "golang/test2", Score: 7}},
		GradingManualRubrics: rubrics,
		Projects: []*database.AssignmentProjects{
			{
				Assignment: database.Assignment{Name: "Assignment 1"},
				GradingManualResults: []*database.ManualGradingResult{
					{Rubric: *rubrics[0], Score: 8},
				},
				GradingJUnitTestResult: &gradingJUnitTestResult,
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

	err := GenerateCSVReport(&reports, assignments, rubrics, nil, true)
	csvReports := reports.String()
	assert.NoError(t, err)
	assert.NotNil(t, csvReports)

	r := csv.NewReader(strings.NewReader(csvReports))
	r.Comma = ';'
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
	assert.Equal(t, "QualityMaxScore", records[0][7])
	assert.Equal(t, "10", records[1][7])
	assert.Equal(t, "AutogradingScore", records[0][8])
	assert.Equal(t, "2", records[1][8])
	assert.Equal(t, "AutogradingMaxScore", records[0][9])
	assert.Equal(t, "9", records[1][9])
	assert.Equal(t, "MaxScore", records[0][10])
	assert.Equal(t, "19", records[1][10])
	assert.Equal(t, "Score", records[0][11])
	assert.Equal(t, "10", records[1][11])
	assert.Equal(t, "Percentage", records[0][12])
}
