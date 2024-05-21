package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (strl *DefaultController) GetClassroomAssignmentProjects(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()
	assignment := ctx.GetAssignment()

	projects, err := assignmentProjectQuery(c, assignment.ID).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := utils.Map(projects, func(project *database.AssignmentProjects) *ProjectResponse {
		return &ProjectResponse{
			AssignmentProjects: project,
			WebURL:             fmt.Sprintf("/api/v2/classrooms/%s/assignments/%s/projects/%s", classroom.ClassroomID.String(), assignment.ID.String(), project.ID.String()),
		}
	})

	return c.JSON(response)
}
