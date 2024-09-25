package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
)

// @Summary		GetClassroomTeamProjects
// @Description	GetClassroomTeamProjects
// @Id				GetClassroomTeamProjects
// @Tags			project
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Param			teamId		path		string	true	"Team ID"		Format(uuid)
// @Success		200			{array}		api.ProjectResponse
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		403			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/api/v1/classrooms/{classroomId}/teams/{teamId}/projects [get]
func (ctrl *DefaultController) GetClassroomTeamProjects(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()
	team := ctx.GetTeam()

	projects, err := teamProjectQuery(c, team.ID).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := utils.Map(projects, func(project *database.AssignmentProjects) *ProjectResponse {
		return &ProjectResponse{
			AssignmentProjects: project,
			WebURL:             fmt.Sprintf("/api/v1/classrooms/%s/teams/%s/projects/%s/gitlab", classroom.ClassroomID.String(), team.ID.String(), project.ID.String()),
			ReportWebURL:       fmt.Sprintf("/api/v1/classrooms/%s/teams/%s/projects/%s/report/gitlab", classroom.ClassroomID.String(), team.ID.String(), project.ID.String()),
		}
	})

	return c.JSON(response)
}
