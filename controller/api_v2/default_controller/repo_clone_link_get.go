package api

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
)

type RepoCloneLinkResponse struct {
	SshUrlToRepo       string                      `json:"sshUrlToRepo"`
	HttpUrlToRepo      string                      `json:"httpUrlToRepo"`
	AssignmentProjects database.AssignmentProjects `json:"assignmentProjects"`
}

func (ctrl *DefaultController) GetRepoCloneLink(c *fiber.Ctx) (err error) {
	return c.SendStatus(fiber.StatusNotImplemented)
}
