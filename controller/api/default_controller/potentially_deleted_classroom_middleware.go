package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab/model"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) PotentiallyDeletedClassroomMiddleware(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)

	// this requires that this middleware is called after the classroom middleware, so that the userClassroom is available
	classroom := ctx.GetUserClassroom().Classroom

	if !classroom.PotentiallyDeleted {
		return c.Next()
	}

	log := ctx.GetLoggerForHandler("PotentiallyDeletedClassroomMiddleware")
	log = log.With("classroom", classroom)

	repo := ctx.GetGitlabRepository()
	_, err = repo.GetGroupById(classroom.GroupID)
	if err == nil {
		// User has access to the group --> group access token got revoked

		classroom.PotentiallyDeleted = false
		classroom.Archived = true
		err := query.Classroom.WithContext(ctx.Context()).Save(&classroom)
		if err != nil {
			return c.Next()
		}

		log.Info("classroom archived due to revoked group access token", "groupID", classroom.GroupID)
		return c.Next()
	}

	var gitLabError *model.GitLabError
	if !errors.As(err, &gitLabError) {
		return c.Next()
	}

	if gitLabError.Response.StatusCode == 404 {
		_, err := query.Classroom.WithContext(ctx.Context()).Delete(&classroom)
		if err != nil {
			return c.Next()
		}
		log.Info("classroom deleted due to group deletion via GitLab.", "groupID", classroom.GroupID)
		return fiber.NewError(fiber.StatusNotFound, "Classroom got deleted via GitLab.")
	}

	return c.Next()
}
