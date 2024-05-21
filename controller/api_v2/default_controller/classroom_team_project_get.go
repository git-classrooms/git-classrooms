package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetClassroomTeamProject
// @Description	GetClassroomTeamProject
// @Id				GetClassroomTeamProject
// @Tags			project
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Param			teamId		path		string	true	"Team ID"		Format(uuid)
// @Param			projectId	path		string	true	"Project ID"	Format(uuid)
// @Success		200			{object}	api.ProjectResponse
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		403			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/teams/{teamId}/projects/{projectId} [get]
func (ctrl *DefaultController) GetClassroomTeamProject(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()
	team := ctx.GetTeam()
	project := ctx.GetAssignmentProject()

	response := &ProjectResponse{
		AssignmentProjects: project,
		WebURL:             fmt.Sprintf("/api/v2/classrooms/%s/teams/%s/projects/%s", classroom.ClassroomID.String(), team.ID.String(), project.ID.String()),
	}

	return c.JSON(response)
}
