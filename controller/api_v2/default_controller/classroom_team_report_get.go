package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gitlab.hs-flensburg.de/gitlab-classroom/model/database/query"
	"gitlab.hs-flensburg.de/gitlab-classroom/utils"
	"gitlab.hs-flensburg.de/gitlab-classroom/wrapper/context"
	"time"
)

// @Summary		GetClassroomTeamReport
// @Description	GetClassroomTeamReport
// @Id				GetClassroomTeamReport
// @Tags			report
// @Produce		application/zip
// @Produce		json
// @Param			classroomId	path		string	true	"Classroom ID"	Format(uuid)
// @Param			teamId		path		string	true	"Team ID"		Format(uuid)
// @Success		200			{file}		application/zip
// @Success		200			{array}		utils.ReportDataItem
// @Failure		400			{object}	HTTPError
// @Failure		401			{object}	HTTPError
// @Failure		404			{object}	HTTPError
// @Failure		500			{object}	HTTPError
// @Router			/api/v2/classrooms/{classroomId}/teams/{teamId}/grading/report [get]
func (ctrl *DefaultController) GetClassroomTeamReport(c *fiber.Ctx) (err error) {
	ctx := context.Get(c)
	classroom := ctx.GetUserClassroom()
	team := ctx.GetTeam()

	assignments, err := assignmentGradingQuery(c, classroom.ClassroomID).
		Join(query.AssignmentProjects, query.AssignmentProjects.AssignmentID.EqCol(query.Assignment.ID)).
		Where(query.AssignmentProjects.TeamID.Eq(team.ID)).
		Find()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else if len(assignments) == 0 {
		return fiber.NewError(fiber.StatusNotFound, "No Assignments found")
	}

	acceptHeader := c.Get("Accept")
	if acceptHeader == "application/json" {
		jsonReports, err := utils.GenerateReports(assignments, &team.ID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(jsonReports)
	}

	c.Set(fiber.HeaderContentType, "application/zip")
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=report_%s_%s_%s.zip", time.Now().Format(time.DateOnly), classroom.Classroom.Name, team.Name))
	return utils.GenerateCSVReports(c.Response().BodyWriter(), assignments, &team.ID)
}
