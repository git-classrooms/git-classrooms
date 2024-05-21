package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) GetClassroomProject(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	project := ctx.GetAssignmentProject()

	response := &ProjectResponse{
		AssignmentProjects: project,
		WebURL:             fmt.Sprintf("/api/v2/classrooms/%s/projects/%s", ctx.GetUserClassroom().ClassroomID, project.ID.String()),
	}

	return c.JSON(response)
}
