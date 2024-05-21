package api

import (
	"github.com/gofiber/fiber/v2"
	apiV2 "gitlab.hs-flensburg.de/gitlab-classroom/controller/api_v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
)

func (ctrl *DefaultController) ViewableClassroomMiddleware() fiber.Handler {
	var validateViewable apiV2.ValidateUserFunc = func(classroom database.UserClassrooms) bool {
		return classroom.Classroom.StudentsViewAllProjects || classroom.Role != database.Student
	}

	return ctrl.ValidateUserMiddleware(validateViewable)
}
