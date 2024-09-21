package utils

import (
	"log"

	dbModel "gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gorm.io/gorm"
)

// MigrateDatabase migrates the database schema to the latest version.
func MigrateDatabase(db *gorm.DB) error {
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	log.Println("Running database migrations")
	return db.AutoMigrate(
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
}
