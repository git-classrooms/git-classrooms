package default_controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

type getOwnedClassroomTeamProjectResponse struct {
	database.AssignmentProjects
	ProjectPath string `json:"projectPath"`
}

//	@Summary		Get all Projects of current team
//	@Description	Get all gitlab projects of the current team
//	@Id				GetOwnedClassroomTeamProjects
//	@Tags			team, project
//	@Accept			json
//	@Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
//	@Param			teamId		path		string	true	"Team ID"		Format(uuid)
//	@Success		200			{array}		default_controller.getOwnedClassroomTeamProjectResponse
//	@Failure		400			{object}	httputil.HTTPError
//	@Failure		401			{object}	httputil.HTTPError
//	@Failure		404			{object}	httputil.HTTPError
//	@Failure		500			{object}	httputil.HTTPError
//	@Router			/classrooms/owned/{classroomId}/teams/{teamId}/projects [get]
func (ctrl *DefaultController) GetOwnedClassroomTeamProjects(c *fiber.Ctx) error {
	ctx := context.Get(c)

	team := ctx.GetOwnedClassroomTeam()

	responses := make([]getOwnedClassroomTeamProjectResponse, len(team.AssignmentProjects))
	for i, project := range team.AssignmentProjects {
		responses[i] = getOwnedClassroomTeamProjectResponse{
			AssignmentProjects: *project,
			ProjectPath:        fmt.Sprintf("/api/v1/classrooms/owned/%s/assignments/%s/projects/%s/gitlab", team.ClassroomID.String(), project.AssignmentID.String(), project.ID.String()),
		}
	}
	return ctx.JSON(responses)
}
