package utils

import (
	dbModel "gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gorm.io/gorm"
	"log"
)

func MigrateDatabase(db *gorm.DB) error {
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)

	log.Println("Running database migrations")
	return db.AutoMigrate(
		&dbModel.User{},
		&dbModel.Classroom{},
		&dbModel.UserClassrooms{},
		&dbModel.Assignment{},
		&dbModel.AssignmentProjects{},
		&dbModel.ClassroomInvitation{},
		&dbModel.AssignmentInvitation{},
	)
}
