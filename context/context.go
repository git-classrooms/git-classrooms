package context

import (
	"backend/api/repository/go_gitlab_repo"
	"github.com/gofiber/fiber/v2"
)

const (
	gitlabRepoKey = "gitlab-repo"
)

func GetGitlabRepository(c *fiber.Ctx) *go_gitlab_repo.GoGitlabRepo {
	return c.Locals(gitlabRepoKey).(*go_gitlab_repo.GoGitlabRepo)
}

func SetGitlabRepository(c *fiber.Ctx, repo *go_gitlab_repo.GoGitlabRepo) {
	c.Locals(gitlabRepoKey, repo)
}
