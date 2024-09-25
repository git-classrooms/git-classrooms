package api

import (
	"slices"

	"github.com/gofiber/fiber/v2"

	apiController "gitlab.hs-flensburg.de/gitlab-classroom/controller/api"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
)

func (ctrl *DefaultController) RoleMiddleware(validRoles ...database.Role) fiber.Handler {
	var validateRoles apiController.ValidateUserFunc = func(user database.UserClassrooms) bool {
		return slices.Contains(validRoles, user.Role)
	}

	return ctrl.ValidateUserMiddleware(validateRoles)
}
