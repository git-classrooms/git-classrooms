package api

import (
	"github.com/gofiber/fiber/v2"
	apiController "gitlab.hs-flensburg.de/gitlab-classroom/controller/api"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
)

func (ctrl *DefaultController) ViewableClassroomMiddleware() fiber.Handler {
	var validateViewable apiController.ValidateUserFunc = func(classroom database.UserClassrooms) bool {
		return classroom.Classroom.StudentsViewAllProjects || classroom.Role != database.Student
	}

	return ctrl.ValidateUserMiddleware(validateViewable)
}
