package main

import (
	"github.com/google/uuid"
	dbModel "gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gorm.io/gen"
)

type Querier interface {
	// SELECT * FROM @@table
	//   {{where}}
	//     classroom_id = @classroomID
	//     {{if len(teamIDs) != 0}}
	//       AND id NOT IN @teamIDs
	//     {{end}}
	//   {{end}}
	FindByClassroomIDAndNotInTeamIDs(classroomID uuid.UUID, teamIDs ...uuid.UUID) ([]*gen.T, error)
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
		// 	&dbModel.Team{},
		&dbModel.UserClassrooms{},
		&dbModel.Assignment{},
		&dbModel.AssignmentProjects{},
		&dbModel.ClassroomInvitation{},
	)

	g.ApplyInterface(func(Querier) {}, dbModel.Team{})

	g.Execute()
}
