package main

import (
	"github.com/google/uuid"
	dbModel "gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gorm.io/gen"
)

type TeamQuerier interface {
	// SELECT * FROM @@table
	//   {{where}}
	//     classroom_id = @classroomID
	//     {{if len(teamIDs) != 0}}
	//       AND id NOT IN @teamIDs
	//     {{end}}
	//   {{end}}
	FindByClassroomIDAndNotInTeamIDs(classroomID uuid.UUID, teamIDs ...uuid.UUID) ([]*gen.T, error)
}

type ManualGradingRubricQuerier interface {
	// SELECT * FROM @@table
	//   {{where}}
	//     assignment_id = @assignmentID
	//     {{if len(rubricIDs) != 0}}
	//       AND id IN @rubricIDs
	//     {{else}}
	//		 AND 1 = 0
	//     {{end}}
	//   {{end}}
	FindByAssignmentIDAndInRubricIDs(assignmentID uuid.UUID, rubricIDs ...uuid.UUID) ([]*gen.T, error)

	// DELETE FROM @@table
	//   {{where}}
	//     assignment_id = @assignmentID
	//     {{if len(rubricIDs) != 0}}
	//       AND id NOT IN @rubricIDs
	//     {{end}}
	//   {{end}}
	DeleteByAssignmentIDAndNotInRubricIDs(assignmentID uuid.UUID, rubricIDs ...uuid.UUID) error
}

func main() {
	g := gen.NewGenerator(gen.Config{
		WithUnitTest: false,
		OutPath:      "model/database/query",
		Mode:         gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.ApplyBasic(
		&dbModel.User{},
		&dbModel.Classroom{},
		&dbModel.Team{},
		&dbModel.UserClassrooms{},
		&dbModel.Assignment{},
		&dbModel.AssignmentProjects{},
		&dbModel.ClassroomInvitation{},
		&dbModel.ManualGradingRubric{},
		&dbModel.ManualGradingResult{},
		&dbModel.AssignmentJunitTest{},
	)

	g.ApplyInterface(func(TeamQuerier) {}, dbModel.Team{})
	g.ApplyInterface(func(ManualGradingRubricQuerier) {}, dbModel.ManualGradingRubric{})

	g.Execute()
}
