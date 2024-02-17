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

type FiberContext struct {
	*fiber.Ctx
}

func Get(c *fiber.Ctx) *FiberContext {
	return &FiberContext{Ctx: c}
}

func (c *FiberContext) GetGitlabRepository() gitlab.Repository {
	value, ok := c.Locals(gitlabRepoKey).(gitlab.Repository)
	if !ok {
		return nil
	}
	return value
}

func (c *FiberContext) SetGitlabRepository(repo gitlab.Repository) {
	c.Locals(gitlabRepoKey, repo)
}

func (c *FiberContext) SetClassroom(classroom *database.UserClassrooms) {
	c.Locals(classroomKey, classroom)
}

func (c *FiberContext) GetClassroom() *database.UserClassrooms {
	value, ok := c.Locals(classroomKey).(*database.UserClassrooms)
	if !ok {
		return nil
	}
	return value
}
