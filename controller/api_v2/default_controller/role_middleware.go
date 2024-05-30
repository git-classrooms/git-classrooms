package api

import (
	"slices"

	"github.com/gofiber/fiber/v2"
	apiV2 "gitlab.hs-flensburg.de/gitlab-classroom/controller/api_v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
)

func (ctrl *DefaultController) RoleMiddleware(validRoles ...database.Role) fiber.Handler {
	var validateRoles apiV2.ValidateUserFunc = func(user database.UserClassrooms) bool {
		return slices.Contains(validRoles, user.Role)
	}

	return ctrl.ValidateUserMiddleware(validateRoles)
}
