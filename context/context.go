package context

import (
	"de.hs-flensburg.gitlab/gitlab-classroom/api/repository/go_gitlab_repo"
	"github.com/gofiber/fiber/v2"
)

const (
	gitlabRepoKey = "gitlab-repo"
)

func GetGitlabRepository(c *fiber.Ctx) *go_gitlab_repo.GoGitlabRepo {
	value, ok := c.Locals(gitlabRepoKey).(*go_gitlab_repo.GoGitlabRepo)
	if !ok {
		return nil
	}
	return value
}

func SetGitlabRepository(c *fiber.Ctx, repo *go_gitlab_repo.GoGitlabRepo) {
	c.Locals(gitlabRepoKey, repo)
}
