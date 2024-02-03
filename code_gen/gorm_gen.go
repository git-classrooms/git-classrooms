package main

import (
	dbModel "gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gorm.io/gen"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		WithUnitTest: false,
		OutPath:      "model/database/query",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.ApplyBasic(
		&dbModel.User{},
		&dbModel.Classroom{},
		&dbModel.UserClassrooms{},
		&dbModel.Assignment{},
		&dbModel.AssignmentProjects{},
		&dbModel.ClassroomInvitation{},
	)

	g.Execute()
}
