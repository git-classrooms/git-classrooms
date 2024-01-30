package context

import (
	"de.hs-flensburg.gitlab/gitlab-classroom/repository/gitlab"
	"github.com/gofiber/fiber/v2"
)

const (
	gitlabRepoKey = "gitlab-repo"
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
