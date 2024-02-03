package context

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/repository/gitlab"
)

const (
	gitlabRepoKey = "gitlab-repo"
	classroomKey  = "classroom"
)

func GetGitlabRepository(c *fiber.Ctx) gitlab.Repository {
	value, ok := c.Locals(gitlabRepoKey).(gitlab.Repository)
	if !ok {
		return nil
	}
	return value
}

func SetGitlabRepository(c *fiber.Ctx, repo gitlab.Repository) {
	c.Locals(gitlabRepoKey, repo)
}

func SetClassroom(c *fiber.Ctx, classroom *database.UserClassrooms) {
	c.Locals(classroomKey, classroom)
}

func GetClassroom(c *fiber.Ctx) *database.UserClassrooms {
	value, ok := c.Locals(classroomKey).(*database.UserClassrooms)
	if !ok {
		return nil
	}
	return value
}
