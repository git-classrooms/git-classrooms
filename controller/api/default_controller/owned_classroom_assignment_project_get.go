package default_controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type getOwnedClassroomAssignmentProjectResponse struct {
	database.AssignmentProjects
	ProjectPath string `json:"projectPath"`
}

func (ctrl *DefaultController) GetOwnedClassroomAssignmentProject(c *fiber.Ctx) error {
	ctx := context.Get(c)
	classroom := ctx.GetOwnedClassroom()
	project := ctx.GetOwnedClassroomAssignmentProject()
	response := getOwnedClassroomAssignmentProjectResponse{
		AssignmentProjects: *project,
		ProjectPath:        fmt.Sprintf("/api/v1/classrooms/owned/%s/assignments/%s/projects/%s/gitlab", classroom.ID.String(), project.AssignmentID.String(), project.ID.String()),
	}
	return c.JSON(response)
}
