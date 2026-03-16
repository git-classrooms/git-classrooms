package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) TeachingGroupGitlabGroupMiddleware(c *fiber.Ctx) error {
	ctx := context.Get(c)

	teachingGroupId := ctx.GetUserClassroom().Classroom.TeachingGroupID

	if teachingGroupId == nil {
		return fiber.ErrNotFound
	}

	ctx.SetGitlabGroupID(*teachingGroupId)

	return c.Next()
}
