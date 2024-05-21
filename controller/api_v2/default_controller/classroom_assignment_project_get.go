package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) GetClassroomAssignmentProject(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()
	assignment := ctx.GetAssignment()
	project := ctx.GetAssignmentProject()

	response := &ProjectResponse{
		AssignmentProjects: project,
		WebURL:             fmt.Sprintf("/api/v2/classrooms/%s/assignments/%s/projects/%s", classroom.ClassroomID.String(), assignment.ID.String(), project.ID.String()),
	}

	return c.JSON(response)
}
