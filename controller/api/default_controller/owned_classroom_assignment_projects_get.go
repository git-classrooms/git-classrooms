package default_controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

func (ctrl *DefaultController) GetOwnedClassroomAssignmentProjects(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()
	assignment := ctx.GetOwnedClassroomAssignment()

	responses := make([]getOwnedClassroomAssignmentProjectResponse, len(assignment.Projects))
	for i, project := range assignment.Projects {
		responses[i] = getOwnedClassroomAssignmentProjectResponse{
			AssignmentProjects: *project,
			ProjectPath:        fmt.Sprintf("/api/v1/classrooms/owned/%s/assignments/%s/projects/%s/gitlab", classroom.ID.String(), project.AssignmentID.String(), project.ID.String()),
		}
	}

	return c.JSON(responses)
}
