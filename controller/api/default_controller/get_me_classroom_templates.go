package default_controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/context"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
)

func (ctrl *DefaultController) GetMeClassroomTemplates(c *fiber.Ctx) error {
	classroom := context.GetClassroom(c)
	if classroom.Role != database.Owner {
		return fiber.NewError(fiber.StatusForbidden, "only the owner can access the templates")
	}

	search := c.Query("search")

	repo := context.GetGitlabRepository(c)
	projects, err := repo.GetAllProjects(search)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(projects)
}
